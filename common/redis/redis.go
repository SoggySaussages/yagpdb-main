package redis

import (
	"os"
	"strings"

	"github.com/mediocregopher/radix/v3"
)

var Prefix = os.Getenv("YAGPDB_REDIS_PREFIX")

type RedisCmdAction struct {
	a radix.CmdAction
}

func (a RedisCmdAction) action() *redisAction {
	return &redisAction{a.a}
}

type redisAction struct {
	a radix.Action
}

func (a redisAction) action() *redisAction {
	return &a
}

type action interface {
	action() *redisAction
}

type Client struct {
	C radix.Client
}

func (c Client) Do(a action) error {
	return c.C.Do(a.action().a)
}

type RedisClient interface {
	Do(a action) error
}

// NewScanner creates a new Scanner instance which will iterate over the redis
// instance's Client using the ScanOpts.
//
// NOTE if Client is a *Cluster this will not work correctly, use the NewScanner
// method on Cluster instead.
func NewScanner(c Client, o radix.ScanOpts) radix.Scanner {
	if o.Pattern != "" {
		o.Pattern = prefixKey(o.Pattern)
	}
	return radix.NewScanner(c.C, o)
}

// Cmd is used to perform a redis command and retrieve a result. It should not
// be passed into Do more than once.
//
// If the receiver value of Cmd is a primitive, a slice/map, or a struct then a
// pointer must be passed in. It may also be an io.Writer, an
// encoding.Text/BinaryUnmarshaler, or a resp.Unmarshaler. See the package docs
// for more on how results are unmarshaled into the receiver.
func Cmd(rcv interface{}, cmd string, args ...string) RedisCmdAction {
	return RedisCmdAction{radix.Cmd(rcv, cmd, argsPrefixedKeys(cmd, args)...)}
}

// FlatCmd is like Cmd, but the arguments can be of almost any type, and FlatCmd
// will automatically flatten them into a single array of strings. Like Cmd, a
// FlatCmd should not be passed into Do more than once.
//
// FlatCmd does _not_ work for commands whose first parameter isn't a key, or
// (generally) for MSET. Use Cmd for those.
//
// FlatCmd supports using a resp.LenReader (an io.Reader with a Len() method) as
// an argument. *bytes.Buffer is an example of a LenReader, and the resp package
// has a NewLenReader function which can wrap an existing io.Reader.
//
// FlatCmd also supports encoding.Text/BinaryMarshalers. It does _not_ currently
// support resp.Marshaler.
//
// The receiver to FlatCmd follows the same rules as for Cmd.
func FlatCmd(rcv interface{}, cmd, key string, args ...interface{}) RedisCmdAction {
	return RedisCmdAction{radix.FlatCmd(rcv, cmd, Prefix+key, args...)}
}

// WithConn is used to perform a set of independent Actions on the same Conn.
//
// key should be a key which one or more of the inner Actions is going to act
// on, or "" if no keys are being acted on or the keys aren't yet known. key is
// generally only necessary when using Cluster.
//
// The callback function is what should actually carry out the inner actions,
// and the error it returns will be passed back up immediately.
//
// NOTE that WithConn only ensures all inner Actions are performed on the same
// Conn, it doesn't make them transactional. Use MULTI/WATCH/EXEC within a
// WithConn for transactions, or use EvalScript.
func WithConn(key string, fn func(radix.Conn) error) redisAction {
	return redisAction{radix.WithConn(prefixKey(key), fn)}
}

// Pipeline returns an Action which first writes multiple commands to a Conn in
// a single write, then reads their responses in a single read. This reduces
// network delay into a single round-trip.
//
// Run will not be called on any of the passed in CmdActions.
//
// NOTE that, while a Pipeline performs all commands on a single Conn, it
// shouldn't be used by itself for MULTI/EXEC transactions, because if there's
// an error it won't discard the incomplete transaction. Use WithConn or
// EvalScript for transactional functionality instead.
func Pipeline(cmds ...RedisCmdAction) redisAction {
	var commands []radix.CmdAction
	for _, c := range cmds {
		commands = append(commands, c.a)
	}

	return redisAction{radix.Pipeline(commands...)}
}

type RedisPool struct {
	p *radix.Pool
}

func (p *RedisPool) NewPool(network, addr string, size int, opts ...radix.PoolOpt) (err error) {
	p.p, err = radix.NewPool(network, addr, size, opts...)
	return
}

func (p RedisPool) Do(a action) error {
	return p.p.Do(a.action().a)
}

func argsPrefixedKeys(cmd string, args []string) []string {
	cmd = strings.ToUpper(cmd)
	if cmd == "BITOP" && len(args) > 1 {
		keys := prefixKeys(args[1:])
		return append(args[:1], keys...)
	} else if cmd == "XINFO" {
		if len(args) < 2 {
			return args
		}
		args[1] = prefixKey(args[1])
		return args
	} else if cmd == "XGROUP" && len(args) > 1 {
		args[1] = prefixKey(args[1])
		return args
	} else if cmd == "XREAD" || cmd == "XREADGROUP" {
		in, out := findStreamsKeys(args)

		for i := in; i < out; i++ {
			args[i] = prefixKey(args[i])
		}

		return args
	} else if noKeyCmds[cmd] || len(args) == 0 {
		return args
	}
	args[0] = prefixKey(args[0])
	return args
}

func prefixKey(key string) string {
	return Prefix + key
}

func prefixKeys(keys []string) (prefixed []string) {
	for _, k := range keys {
		prefixed = append(prefixed, prefixKey(k))
	}
	return
}

func findStreamsKeys(args []string) (int, int) {
	for i, arg := range args {
		if strings.ToUpper(arg) != "STREAMS" {
			continue
		}

		// after STREAMS only stream keys and IDs can be given and since there must be the same number of keys and ids
		// we can just take half of remaining arguments as keys. If the number of IDs does not match the number of
		// keys the command will fail later when send to Redis so no need for us to handle that case.
		ids := len(args[i+1:]) / 2

		return i + 1, len(args) - ids
	}

	return 0, 0
}

var noKeyCmds = map[string]bool{
	"SENTINEL": true,

	"CLUSTER":   true,
	"READONLY":  true,
	"READWRITE": true,
	"ASKING":    true,

	"AUTH":   true,
	"ECHO":   true,
	"PING":   true,
	"QUIT":   true,
	"SELECT": true,
	"SWAPDB": true,

	"KEYS":      true,
	"MIGRATE":   true,
	"OBJECT":    true,
	"RANDOMKEY": true,
	"WAIT":      true,
	"SCAN":      true,

	"EVAL":    true,
	"EVALSHA": true,
	"SCRIPT":  true,

	"BGREWRITEAOF": true,
	"BGSAVE":       true,
	"CLIENT":       true,
	"COMMAND":      true,
	"CONFIG":       true,
	"DBSIZE":       true,
	"DEBUG":        true,
	"FLUSHALL":     true,
	"FLUSHDB":      true,
	"INFO":         true,
	"LASTSAVE":     true,
	"MONITOR":      true,
	"ROLE":         true,
	"SAVE":         true,
	"SHUTDOWN":     true,
	"SLAVEOF":      true,
	"SLOWLOG":      true,
	"SYNC":         true,
	"TIME":         true,

	"DISCARD": true,
	"EXEC":    true,
	"MULTI":   true,
	"UNWATCH": true,
	"WATCH":   true,
}
