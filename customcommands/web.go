package customcommands

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode/utf8"

	"emperror.dev/errors"
	"github.com/botlabs-gg/yagpdb/v2/common"
	"github.com/botlabs-gg/yagpdb/v2/common/cplogs"
	"github.com/botlabs-gg/yagpdb/v2/common/featureflags"
	prfx "github.com/botlabs-gg/yagpdb/v2/common/prefix"
	"github.com/botlabs-gg/yagpdb/v2/common/pubsub"
	yagtemplate "github.com/botlabs-gg/yagpdb/v2/common/templates"
	"github.com/botlabs-gg/yagpdb/v2/customcommands/models"
	"github.com/botlabs-gg/yagpdb/v2/lib/discordgo"
	"github.com/botlabs-gg/yagpdb/v2/premium"
	"github.com/botlabs-gg/yagpdb/v2/web"
	"github.com/mediocregopher/radix/v3"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"goji.io"
	"goji.io/pat"
)

//go:embed assets/customcommands-editcmd.html
var PageHTMLEditCmd string

//go:embed assets/customcommands.html
var PageHTMLMain string

//go:embed assets/customcommands-database.html
var PageHTMLDatabase string

// GroupForm is the form bindings used when creating or updating groups
type GroupForm struct {
	ID                int64
	Name              string  `valid:",100"`
	GitHub            string  `valid:",100"`
	WhitelistChannels []int64 `valid:"channel,true"`
	BlacklistChannels []int64 `valid:"channel,true"`

	WhitelistRoles []int64 `valid:"role,true"`
	BlacklistRoles []int64 `valid:"role,true"`
}

type SearchForm struct {
	Query string
	Type  string
}

var (
	panelLogKeyNewCommand     = cplogs.RegisterActionFormat(&cplogs.ActionFormat{Key: "customcommands_new_command", FormatString: "Created a new custom command: %d"})
	panelLogKeyUpdatedCommand = cplogs.RegisterActionFormat(&cplogs.ActionFormat{Key: "customcommands_updated_command", FormatString: "Updated custom command: %d"})
	panelLogKeyRemovedCommand = cplogs.RegisterActionFormat(&cplogs.ActionFormat{Key: "customcommands_removed_command", FormatString: "Removed custom command: %d"})

	panelLogKeyNewGroup     = cplogs.RegisterActionFormat(&cplogs.ActionFormat{Key: "customcommands_new_group", FormatString: "Created a new custom command group: %s"})
	panelLogKeyUpdatedGroup = cplogs.RegisterActionFormat(&cplogs.ActionFormat{Key: "customcommands_updated_group", FormatString: "Updated custom command group: %s"})
	panelLogKeyRemovedGroup = cplogs.RegisterActionFormat(&cplogs.ActionFormat{Key: "customcommands_removed_group", FormatString: "Removed custom command group: %d"})
)

// InitWeb implements web.Plugin
func (p *Plugin) InitWeb() {
	web.AddHTMLTemplate("customcommands/assets/customcommands.html", PageHTMLMain)
	web.AddHTMLTemplate("customcommands/assets/customcommands-editcmd.html", PageHTMLEditCmd)
	web.AddSidebarItem(web.SidebarCategoryCustomCommands, &web.SidebarItem{
		Name: "Commands",
		URL:  "customcommands",
		Icon: "fas fa-code",
	})

	web.AddHTMLTemplate("customcommands/assets/customcommands-database.html", PageHTMLDatabase)
	web.AddSidebarItem(web.SidebarCategoryCustomCommands, &web.SidebarItem{
		Name: "Database",
		URL:  "customcommands/database",
		Icon: "fas fa-database",
	})

	getHandler := web.ControllerHandler(handleCommands, "cp_custom_commands")
	getCmdHandler := web.ControllerHandler(handleGetCommand, "cp_custom_commands_edit_cmd")
	getGroupHandler := web.ControllerHandler(handleGetCommandsGroup, "cp_custom_commands")
	getDBHandler := web.ControllerHandler(handleGetDatabase, "cp_custom_commands_database")

	subMux := goji.SubMux()
	web.CPMux.Handle(pat.New("/customcommands"), subMux)
	web.CPMux.Handle(pat.New("/customcommands/*"), subMux)
	web.CPMux.Handle(pat.New("/customcommands/database"), subMux)

	subMux.Use(func(inner http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			g, templateData := web.GetBaseCPContextData(r.Context())
			strTriggerTypes := map[int]string{}
			for k, v := range triggerStrings {
				strTriggerTypes[int(k)] = v
			}
			templateData["CCTriggerTypes"] = strTriggerTypes
			templateData["CommandPrefix"], _ = prfx.GetCommandPrefixRedis(g.ID)

			inner.ServeHTTP(w, r)
		}
		return http.HandlerFunc(h)
	})

	subMux.Handle(pat.Get(""), getHandler)
	subMux.Handle(pat.Get("/"), getHandler)

	subMux.Handle(pat.Get("/database"), getDBHandler)
	subMux.Handle(pat.Get("/database/"), getDBHandler)
	subMux.Handle(pat.Post("/database/delete/:id"), web.ControllerPostHandler(handleDeleteDatabaseEntry, getDBHandler, nil))

	subMux.Handle(pat.Get("/commands/:cmd/"), getCmdHandler)

	subMux.Handle(pat.Get("/groups/:group/"), web.ControllerHandler(handleGetCommandsGroup, "cp_custom_commands"))
	subMux.Handle(pat.Get("/groups/:group"), web.ControllerHandler(handleGetCommandsGroup, "cp_custom_commands"))

	newCommandHandler := web.ControllerPostHandler(handleNewCommand, nil, nil)
	subMux.Handle(pat.Post("/commands/new"), newCommandHandler)
	subMux.Handle(pat.Post("/commands/:cmd/update"), web.ControllerPostHandler(handleUpdateCommand, getCmdHandler, CustomCommand{}))
	subMux.Handle(pat.Post("/commands/:cmd/delete"), web.ControllerPostHandler(handleDeleteCommand, getHandler, nil))
	subMux.Handle(pat.Post("/commands/:cmd/run_now"), web.ControllerPostHandler(handleRunCommandNow, getCmdHandler, nil))
	subMux.Handle(pat.Post("/commands/:cmd/update_and_run"), web.ControllerPostHandler(handleUpdateAndRunNow, getCmdHandler, CustomCommand{}))
	subMux.Handle(pat.Post("/commands/:cmd/git_pull"), web.ControllerPostHandler(handleGitPull, getCmdHandler, CustomCommand{}))

	subMux.Handle(pat.Post("/creategroup"), web.ControllerPostHandler(handleNewGroup, getHandler, GroupForm{}))
	subMux.Handle(pat.Post("/groups/:group/update"), web.ControllerPostHandler(handleUpdateGroup, getGroupHandler, GroupForm{}))
	subMux.Handle(pat.Post("/groups/:group/delete"), web.ControllerPostHandler(handleDeleteGroup, getHandler, nil))

	web.ServerPublicMux.Handle(pat.Get("/customcommands/commands/:id"), PublicCommandMW(getCmdHandler))
	web.ServerPublicMux.Handle(pat.Get("/customcommands/commands/:id/"), PublicCommandMW(getCmdHandler))
}

func handleGetDatabase(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	var err error

	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)

	page := 0
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			templateData.AddAlerts(web.ErrorAlert("Failed parsing Page"))
		}
	}

	if page < 1 {
		page = 1
	}

	queryType := r.URL.Query().Get("type")
	query := r.URL.Query().Get("query")
	if len(query) > 0 {
		templateData["Query"] = query
		templateData["QueryType"] = queryType
	}

	result, total, err := getDatabaseEntries(ctx, activeGuild.ID, page, queryType, query, 100)
	if err != nil {
		return templateData, err
	}
	totalPages := int(math.Ceil((float64(total) / 100)))
	if page > totalPages {
		page = totalPages
	}

	entries := convertEntries(result)
	templateData["DBEntries"] = entries
	templateData["TotalPages"] = totalPages
	templateData["Page"] = page

	return templateData, nil
}

func handleDeleteDatabaseEntry(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)

	id, err := strconv.ParseInt(pat.Param(r, "id"), 10, 64)
	if err != nil {
		return templateData, err
	}

	_, err = models.TemplatesUserDatabases(qm.Where("guild_id = ? AND id = ?", activeGuild.ID, id)).DeleteAll(ctx, common.PQ)
	if err != nil {
		return templateData, err
	}

	return templateData.AddAlerts(), nil
}

func handleCommands(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)

	groupID := int64(0)
	if v, ok := templateData["CurrentGroupID"]; ok {
		groupID = v.(int64)
	}

	var langBuiltins strings.Builder
	for k := range yagtemplate.StandardFuncMap {
		langBuiltins.WriteString(" " + k)
	}

	templateData["HLJSBuiltins"] = langBuiltins.String()

	count, err := models.CustomCommands(qm.Where("guild_id = ?", activeGuild.ID)).CountG(ctx)
	if err != nil {
		return templateData, err
	}

	updateTemplateWithCountData(int(count), templateData, ctx)

	return serveGroupSelected(r, templateData, groupID, activeGuild.ID)
}

func handleGetCommand(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	activeGuild, templateData := web.GetBaseCPContextData(r.Context())

	ccID, err := strconv.ParseInt(pat.Param(r, "cmd"), 10, 64)
	if err != nil {
		return templateData, errors.WithStackIf(err)
	}

	cc, err := models.CustomCommands(
		models.CustomCommandWhere.GuildID.EQ(activeGuild.ID),
		models.CustomCommandWhere.LocalID.EQ(ccID)).OneG(r.Context())
	if err != nil {
		return templateData, errors.WithStackIf(err)
	}

	templateData["CC"] = cc
	templateData["Commands"] = true
	templateData["IsGuildPremium"] = premium.ContextPremium(r.Context())
	templateData["GitHubFilepath"] = "<cc not in group>"
	group, err := models.CustomCommandGroups(qm.Where("guild_id = ? AND id = ?", activeGuild.ID, cc.GroupID.Int64)).OneG(r.Context())
	if err == nil {
		if group.GitHub == "" {
			templateData["GitHubFilepath"] = "<github not specified in group>"
		} else {
			templateData["GitHubFilepath"] = fmt.Sprintf("%s/%d/%d.yag", group.GitHub, activeGuild.ID, cc.LocalID)
			templateData["GitHubAvailable"] = true
			templateData["GitHub"] = group.GitHub
		}
	}
	templateData["PublicLink"] = fmt.Sprintf("%s/public/%d/customcommands/commands/%d", web.BaseURL(), activeGuild.ID, cc.LocalID)

	return serveGroupSelected(r, templateData, cc.GroupID.Int64, activeGuild.ID)
}

func handleGetCommandsGroup(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)
	groupID, _ := strconv.ParseInt(pat.Param(r, "group"), 10, 64)

	count, err := models.CustomCommands(qm.Where("guild_id = ?", activeGuild.ID)).CountG(ctx)
	if err != nil {
		return templateData, err
	}

	updateTemplateWithCountData(int(count), templateData, ctx)

	return serveGroupSelected(r, templateData, groupID, activeGuild.ID)
}

func serveGroupSelected(r *http.Request, templateData web.TemplateData, groupID int64, guildID int64) (web.TemplateData, error) {
	templateData["GetCCIntervalType"] = tmplGetCCIntervalTriggerType
	templateData["GetCCInterval"] = tmplGetCCInterval

	_, ok := templateData["CustomCommands"]
	if !ok {
		var err error
		var commands []*models.CustomCommand
		if groupID == 0 {
			commands, err = models.CustomCommands(qm.Where("guild_id = ? AND group_id IS NULL", guildID), qm.OrderBy("local_id asc")).AllG(r.Context())
		} else {
			commands, err = models.CustomCommands(qm.Where("guild_id = ? AND group_id = ?", guildID, groupID), qm.OrderBy("local_id asc")).AllG(r.Context())
		}
		if err != nil {
			return templateData, err
		}

		templateData["CustomCommands"] = commands
	}

	commandsGroups, err := models.CustomCommandGroups(qm.Where("guild_id = ?", guildID), qm.OrderBy("id asc")).AllG(r.Context())
	if err != nil {
		return templateData, err
	}

	for _, v := range commandsGroups {
		if v.ID == groupID {
			templateData["CurrentCommandGroup"] = v
			break
		}
	}

	templateData["CommandGroups"] = commandsGroups

	return templateData, nil
}

func handleNewCommand(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)

	groupID, _ := strconv.ParseInt(r.FormValue("GroupID"), 10, 64)
	if groupID != 0 {
		// make sure we aren't trying to pull any tricks with the group id
		c, err := models.CustomCommandGroups(qm.Where("guild_id = ? AND id = ?", activeGuild.ID, groupID)).CountG(ctx)
		if err != nil {
			return templateData, err
		}

		if c < 1 {
			return templateData.AddAlerts(web.ErrorAlert("Unknown group")), nil
		}

		templateData["CurrentGroupID"] = groupID
	}

	c, err := models.CustomCommands(qm.Where("guild_id = ?", activeGuild.ID)).CountG(ctx)
	if err != nil {
		return templateData, err
	}

	if int(c) >= MaxCommandsForContext(ctx) {
		return templateData, web.NewPublicError(fmt.Sprintf("Max %d custom commands allowed (or %d for premium servers)", MaxCommands, MaxCommandsPremium))
	}

	localID, err := common.GenLocalIncrID(activeGuild.ID, "custom_command")
	if err != nil {
		return templateData, errors.WrapIf(err, "error generating local id")
	}

	dbModel := &models.CustomCommand{
		GuildID: activeGuild.ID,
		LocalID: localID,

		Disabled:   false,
		ShowErrors: true,

		TimeTriggerExcludingDays:  []int64{},
		TimeTriggerExcludingHours: []int64{},

		Responses: []string{"Edit this to change the output of the custom command {{.CCID}}!"},
	}

	if groupID != 0 {
		dbModel.GroupID = null.Int64From(groupID)
	}

	err = dbModel.InsertG(ctx, boil.Infer())
	if err != nil {
		return templateData, err
	}

	featureflags.MarkGuildDirty(activeGuild.ID)

	http.Redirect(w, r, fmt.Sprintf("/manage/%d/customcommands/commands/%d/", activeGuild.ID, localID), http.StatusSeeOther)

	go cplogs.RetryAddEntry(web.NewLogEntryFromContext(r.Context(), panelLogKeyNewCommand, &cplogs.Param{Type: cplogs.ParamTypeInt, Value: dbModel.LocalID}))

	pubsub.EvictCacheSet(cachedCommandsMessage, activeGuild.ID)
	return templateData, nil
}

func handleUpdateCommand(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)

	cmdEdit := ctx.Value(common.ContextKeyParsedForm).(*CustomCommand)
	cmdSaved, err := models.FindCustomCommandG(context.Background(), activeGuild.ID, int64(cmdEdit.ID))
	if cmdSaved.Disabled == true && cmdEdit.ToDBModel().Disabled == false {
		c, err := models.CustomCommands(qm.Where("guild_id = ? and disabled = false", activeGuild.ID)).CountG(ctx)
		if err != nil {
			return templateData, err
		}
		if int(c) >= MaxCommandsForContext(ctx) {
			return templateData, web.NewPublicError(fmt.Sprintf("Max %d enabled custom commands allowed (or %d for premium servers)", MaxCommands, MaxCommandsPremium))
		}
	}

	// ensure that the group specified is owned by this guild
	if cmdEdit.GroupID != 0 {
		c, err := models.CustomCommandGroups(qm.Where("guild_id = ? AND id = ?", activeGuild.ID, cmdEdit.GroupID)).CountG(ctx)
		if err != nil {
			return templateData, err
		}

		if c < 1 {
			return templateData.AddAlerts(web.ErrorAlert("Unknown group")), nil
		}
	}

	if !premium.ContextPremium(ctx) && cmdEdit.TriggerOnEdit {
		return templateData.AddAlerts(web.ErrorAlert("`Trigger on edits` is a premium feature, your command wasn't saved, please save again after disabling `Trigger on edits`")), nil
	}

	dbModel := cmdEdit.ToDBModel()

	templateData["CurrentGroupID"] = dbModel.GroupID.Int64

	dbModel.GuildID = activeGuild.ID
	dbModel.LocalID = cmdEdit.ID
	dbModel.TriggerType = int(triggerTypeFromForm(cmdEdit.TriggerTypeForm))
	// check low interval limits
	if dbModel.TriggerType == int(CommandTriggerInterval) && dbModel.TimeTriggerInterval <= 10 {
		if dbModel.TimeTriggerInterval < 5 {
			dbModel.TimeTriggerInterval = 5
		}

		ok, err := checkIntervalLimits(ctx, activeGuild.ID, dbModel.LocalID, templateData)
		if err != nil || !ok {
			return templateData, err
		}
	}

	_, err = dbModel.UpdateG(ctx, boil.Blacklist("last_run", "next_run", "local_id", "guild_id", "last_error", "last_error_time", "run_count"))
	if err != nil {
		return templateData, nil
	}

	// create, update or remove the next run time and scheduled event
	if dbModel.TriggerType == int(CommandTriggerInterval) {
		// need the last run time
		fullModel, err := models.CustomCommands(qm.Where("guild_id = ? AND local_id = ?", activeGuild.ID, dbModel.LocalID)).OneG(ctx)
		if err != nil {
			web.CtxLogger(ctx).WithError(err).Error("failed retrieving full model")
		} else {
			err = UpdateCommandNextRunTime(fullModel, true, true)
		}
	} else {
		err = DelNextRunEvent(activeGuild.ID, dbModel.LocalID)
	}

	if err != nil {
		web.CtxLogger(ctx).WithError(err).WithField("guild", dbModel.GuildID).Error("failed updating next custom command run time")
	}

	go cplogs.RetryAddEntry(web.NewLogEntryFromContext(r.Context(), panelLogKeyUpdatedCommand, &cplogs.Param{Type: cplogs.ParamTypeInt, Value: dbModel.LocalID}))

	pubsub.EvictCacheSet(cachedCommandsMessage, activeGuild.ID)
	return templateData, err
}

func handleDeleteCommand(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)

	cmdID, err := strconv.ParseInt(pat.Param(r, "cmd"), 10, 64)
	if err != nil {
		return templateData, err
	}

	cmd, err := models.CustomCommands(qm.Where("guild_id = ? AND local_id = ?", activeGuild.ID, cmdID)).OneG(ctx)
	if err != nil {
		return templateData, err
	}

	groupID := cmd.GroupID.Int64
	if groupID != 0 {
		templateData["CurrentGroupID"] = groupID
	}

	_, err = cmd.DeleteG(ctx)
	if err != nil {
		return templateData, err
	}

	go cplogs.RetryAddEntry(web.NewLogEntryFromContext(r.Context(), panelLogKeyRemovedCommand, &cplogs.Param{Type: cplogs.ParamTypeInt, Value: cmd.LocalID}))

	err = DelNextRunEvent(cmd.GuildID, cmd.LocalID)
	featureflags.MarkGuildDirty(activeGuild.ID)
	pubsub.EvictCacheSet(cachedCommandsMessage, activeGuild.ID)
	return templateData, err
}

const RunCmdCooldownSeconds = 5

func keyRunCmdCooldown(guildID, userID int64) string {
	return "custom_command_run_now_cooldown:" + discordgo.StrID(guildID) + ":" + discordgo.StrID(userID)
}

func checkSetCooldown(guildID, userID int64) (bool, error) {
	var resp string
	err := common.RedisPool.Do(radix.FlatCmd(&resp, "SET", keyRunCmdCooldown(guildID, userID), true, "EX", RunCmdCooldownSeconds, "NX"))
	if err != nil {
		return false, err
	}

	return resp == "OK", nil
}

func handleRunCommandNow(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)
	member := web.ContextMember(ctx)
	if member == nil {
		templateData.AddAlerts()
		return templateData, nil
	}

	cmdID, err := strconv.ParseInt(pat.Param(r, "cmd"), 10, 64)
	if err != nil {
		return templateData, err
	}

	// only interval commands can be ran from the dashboard currently
	cmd, err := models.CustomCommands(qm.Where("guild_id = ? AND local_id = ? AND trigger_type = 5", activeGuild.ID, cmdID)).OneG(context.Background())
	if err != nil {
		return templateData, err
	}

	if cmd.Disabled {
		templateData.AddAlerts(web.ErrorAlert("This command is disabled, cannot run a disabled command"))
		return templateData, nil
	}

	ok, err := checkSetCooldown(activeGuild.ID, member.User.ID)
	if err != nil {
		return templateData, err
	}

	if !ok {
		templateData.AddAlerts(web.ErrorAlert("You're on cooldown, wait before trying again"))
		return templateData, nil
	}

	go pubsub.Publish("custom_commands_run_now", activeGuild.ID, cmd)

	return templateData, nil
}

func handleUpdateAndRunNow(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	updateData, err := handleUpdateCommand(w, r)
	if err != nil {
		return updateData, err
	}
	return handleRunCommandNow(w, r)
}

func handleGitPull(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	updateData, err := handleUpdateCommand(w, r)
	if err != nil {
		return updateData, err
	}
	ctx := r.Context()
	activeGuild, _ := web.GetBaseCPContextData(ctx)
	cmd := ctx.Value(common.ContextKeyParsedForm).(*CustomCommand)
	model, err := models.CustomCommandGroups(qm.Where("guild_id = ? AND id = ?", activeGuild.ID, cmd.GroupID)).OneG(ctx)
	if err != nil {
		return updateData, err
	}
	go cloneCCRepo(activeGuild.ID, cmd.GroupID, model.GitHub)
	return updateData, err
}

// allow for max 5 triggers with intervals of less than 10 minutes
func checkIntervalLimits(ctx context.Context, guildID int64, cmdID int64, templateData web.TemplateData) (ok bool, err error) {
	num, err := models.CustomCommands(qm.Where("guild_id = ? AND local_id != ? AND trigger_type = 5 AND time_trigger_interval <= 10", guildID, cmdID)).CountG(ctx)
	if err != nil {
		return false, err
	}

	if num < 5 {
		return true, nil
	}

	templateData.AddAlerts(web.ErrorAlert("You can have max 5 triggers on less than 10 minute intervals"))
	return false, nil
}

func handleNewGroup(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)

	newGroup := ctx.Value(common.ContextKeyParsedForm).(*GroupForm)

	numCurrentGroups, err := models.CustomCommandGroups(qm.Where("guild_id = ?", activeGuild.ID)).CountG(ctx)
	if err != nil {
		return templateData, err
	}

	if numCurrentGroups >= MaxGroups {
		return templateData, web.NewPublicError(fmt.Sprintf("Max %d custom command groups", MaxGroups))
	}

	dbModel := &models.CustomCommandGroup{
		GuildID: activeGuild.ID,
		Name:    newGroup.Name,
	}

	err = dbModel.InsertG(ctx, boil.Infer())
	if err != nil {
		return templateData, err
	}

	go cplogs.RetryAddEntry(web.NewLogEntryFromContext(r.Context(), panelLogKeyNewGroup, &cplogs.Param{Type: cplogs.ParamTypeString, Value: newGroup.Name}))

	templateData["CurrentGroupID"] = dbModel.ID

	pubsub.EvictCacheSet(cachedCommandsMessage, activeGuild.ID)
	return templateData, nil
}

func handleUpdateGroup(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)

	groupForm := ctx.Value(common.ContextKeyParsedForm).(*GroupForm)

	id, _ := strconv.ParseInt(pat.Param(r, "group"), 10, 64)
	model, err := models.CustomCommandGroups(qm.Where("guild_id = ? AND id = ?", activeGuild.ID, id)).OneG(ctx)
	if err != nil {
		return templateData, err
	}

	model.WhitelistChannels = groupForm.WhitelistChannels
	model.IgnoreChannels = groupForm.BlacklistChannels
	model.WhitelistRoles = groupForm.WhitelistRoles
	model.IgnoreRoles = groupForm.BlacklistRoles
	model.Name = groupForm.Name
	model.GitHub = groupForm.GitHub

	go cloneCCRepo(activeGuild.ID, model.ID, model.GitHub)

	_, err = model.UpdateG(ctx, boil.Infer())
	if err == nil {
		go cplogs.RetryAddEntry(web.NewLogEntryFromContext(r.Context(), panelLogKeyUpdatedGroup, &cplogs.Param{Type: cplogs.ParamTypeString, Value: model.Name}))
	}

	pubsub.EvictCacheSet(cachedCommandsMessage, activeGuild.ID)
	return templateData, err
}

func cloneCCRepo(guildID, modelID int64, rawPath string) {
	delDir(fmt.Sprintf("cc-github/%d-%d", guildID, modelID))
	gitArgs := strings.Split(strings.Split(rawPath, "//")[1], "/")
	gitArgs = append([]string{"https:/"}, gitArgs...)
	repo := strings.Join(gitArgs, "/")
	var subDir string
	if len(gitArgs) > 4 && gitArgs[4] != "" {
		repo = strings.Join(gitArgs[:4], "/")
		subDir = strings.Join(gitArgs[4:], "/")
	}
	cmd := exec.Command("git", "clone", repo, fmt.Sprintf("%d-%d", guildID, modelID))
	if subDir != "" {
		cmd = exec.Command("git", "clone", repo, "temp")
	}
	cmd.Dir = "cc-github"
	if subDir == "" {
		runCmdLogErr(cmd)
	} else {
		runSubDir(cmd, subDir, guildID, modelID)
	}
}

func runSubDir(cmd *exec.Cmd, subDir string, guildID, modelID int64) {
	ok := runCmdLogErr(cmd)
	if ok {
		cmd = exec.Command("cp", "-rv", "temp/"+subDir, fmt.Sprintf("%d-%d", guildID, modelID))
		cmd.Dir = "cc-github"
		runCmdLogErr(cmd)
		delDir("cc-github/temp")
	}
}

func delDir(path string) {
	dr, err := os.Open(path)
	if err == nil {
		defer dr.Close()
		logger.Infof("attempting to delete directory %s...", dr.Name())
		files, err := dr.Readdir(-1)
		if err == nil {
			for _, file := range files {
				if file.IsDir() {
					delDir(path + "/" + file.Name())
				}
				err = os.RemoveAll(path + "/" + file.Name())
				if err != nil {
					continue
				}
			}
		}
		err = os.Remove(path)
		if err != nil {
			logger.WithError(err)
		} else {
			logger.Infof("completed deleting directory %s.", dr.Name())
		}
	}
}

func handleDeleteGroup(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)

	id, err := strconv.ParseInt(pat.Param(r, "group"), 10, 64)
	if err != nil {
		return templateData, err
	}

	rows, err := models.CustomCommandGroups(qm.Where("guild_id = ? AND id = ?", activeGuild.ID, id)).DeleteAll(ctx, common.PQ)
	if err != nil {
		return templateData, err
	}

	if rows > 0 {
		go cplogs.RetryAddEntry(web.NewLogEntryFromContext(r.Context(), panelLogKeyRemovedGroup, &cplogs.Param{Type: cplogs.ParamTypeInt, Value: id}))
	}

	pubsub.EvictCacheSet(cachedCommandsMessage, activeGuild.ID)
	return templateData, err
}

func triggerTypeFromForm(str string) CommandTriggerType {
	switch str {
	case "none":
		return CommandTriggerNone
	case "prefix":
		return CommandTriggerStartsWith
	case "regex":
		return CommandTriggerRegex
	case "contains":
		return CommandTriggerContains
	case "exact":
		return CommandTriggerExact
	case "command":
		return CommandTriggerCommand
	case "reaction":
		return CommandTriggerReaction
	case "interval_minutes", "interval_hours":
		return CommandTriggerInterval
	case "component":
		return CommandTriggerComponent
	case "modal":
		return CommandTriggerModal
	default:
		return CommandTriggerCommand

	}
}

func CheckLimits(in ...string) bool {
	for _, v := range in {
		if utf8.RuneCountInString(v) > 2000 {
			return false
		}
	}
	return true
}

// returns 1 for hours, 0 for minutes, -1 otherwise
func tmplGetCCIntervalTriggerType(cc *models.CustomCommand) int {
	if cc.TriggerType != int(CommandTriggerInterval) {
		return -1
	}

	if (cc.TimeTriggerInterval % 60) == 0 {
		return 1
	}

	return 0
}

// returns the proper interval number dispalyed, depending on if it can be rounded to hours or not
func tmplGetCCInterval(cc *models.CustomCommand) int {
	if tmplGetCCIntervalTriggerType(cc) == 1 {
		return cc.TimeTriggerInterval / 60
	}

	return cc.TimeTriggerInterval
}

var _ web.PluginWithServerHomeWidget = (*Plugin)(nil)

func (p *Plugin) LoadServerHomeWidget(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ag, templateData := web.GetBaseCPContextData(r.Context())

	templateData["WidgetTitle"] = "Custom Commands"
	templateData["SettingsPath"] = "/customcommands"

	numCustomCommands, err := models.CustomCommands(qm.Where("guild_id = ?", ag.ID)).CountG(r.Context())

	format := `<p>Number of custom commands: <code>%d</code></p>`

	templateData["WidgetBody"] = template.HTML(fmt.Sprintf(format, numCustomCommands))

	if numCustomCommands > 0 {
		templateData["WidgetEnabled"] = true
	} else {
		templateData["WidgetDisabled"] = true
	}

	return templateData, err
}

func updateTemplateWithCountData(count int, templateData web.TemplateData, ctx context.Context) {
	maxCommands := MaxCommandsForContext(ctx)
	templateData["CCCount"] = count
	templateData["CCLimit"] = maxCommands

	additionalMessage := ""
	if premium.ContextPremiumTier(ctx) != premium.PremiumTierPremium {
		additionalMessage = fmt.Sprintf("(You may increase the limit upto %d with YAGPDB premium)", MaxCommandsPremium)
	}
	templateData["AdditionalMessage"] = additionalMessage
}

func PublicCommandMW(inner http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		defer func() { inner.ServeHTTP(w, r) }()

		ctx := r.Context()
		activeGuild, templateData := web.GetBaseCPContextData(ctx)
		cc := templateData["CC"].(*models.CustomCommand)

		if read, _ := web.IsAdminRequest(ctx, r); read {
			http.Redirect(w, r, fmt.Sprintf("/manage/%d/customcommands/commands/%d/", activeGuild.ID, cc.LocalID), http.StatusSeeOther)
		} else {
			templateData["PublicView"] = true
			if !cc.Public {
				templateData["CC"] = &models.CustomCommand{
					GuildID: activeGuild.ID,
					LocalID: cc.LocalID,

					Disabled:   false,
					ShowErrors: true,

					TimeTriggerExcludingDays:  []int64{},
					TimeTriggerExcludingHours: []int64{},

					Responses: []string{"Edit this to change the output of the custom command {{.CCID}}!"},
				}
				templateData = templateData.AddAlerts(web.ErrorAlert("This command has not been made public by a server admin."))
			}
		}
		r = r.WithContext(context.WithValue(ctx, common.ContextKeyTemplateData, templateData))
	}

	return http.HandlerFunc(mw)
}
