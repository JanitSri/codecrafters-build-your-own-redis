package parser

import (
	"bytes"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/JanitSri/codecrafters-build-your-own-redis/data"
)

var ErrInvalidNumberOfArguments = errors.New("invalid number of arguments")
var ErrInvalidArgument = errors.New("invalid argument")
var ErrInvalidCharacterError = errors.New("invalid error type")

func writeBulkString(s string) []byte {
	l := strconv.Itoa(len(s))
	var b bytes.Buffer
	b.WriteString(BULK_STRING)
	b.WriteString(l)
	b.WriteString(REDIS_TERMINATOR)
	b.WriteString(s)
	b.WriteString(REDIS_TERMINATOR)
	return b.Bytes()
}

type Flag struct {
	name  string
	value string
}

func NewFlag(name, value string) *Flag {
	return &Flag{
		name,
		value,
	}
}

type BaseCommand struct {
	args  []string
	flags []*Flag
}

type Command interface {
	Execute(data.DataStore) []byte
}

type PingCommand struct {
}

func NewPingCommand() *PingCommand {
	return &PingCommand{}
}

func (pc *PingCommand) Execute(ds data.DataStore) []byte {
	log.Println("ponging...")

	var b bytes.Buffer
	b.WriteString(PONG)
	return b.Bytes()
}

type EchoCommand struct {
	BaseCommand
}

func NewEchoCommand(args []string) *EchoCommand {
	return &EchoCommand{
		BaseCommand{
			args,
			nil,
		},
	}
}

func (ec *EchoCommand) Execute(ds data.DataStore) []byte {
	log.Println("echoing...")

	if len(ec.args) != 1 {
		log.Fatal(ErrInvalidNumberOfArguments)
	}

	return writeBulkString(ec.args[0])
}

type SetCommand struct {
	BaseCommand
}

func NewSetCommand(args []string, flags []*Flag) *SetCommand {
	return &SetCommand{
		BaseCommand{
			args,
			flags,
		},
	}
}

func (sc *SetCommand) Execute(ds data.DataStore) []byte {
	log.Println("setting...")

	if len(sc.args) != 2 {
		log.Fatal(ErrInvalidNumberOfArguments)
	}

	v := data.NewRedisValue(sc.args[1], time.Time{})
	ds.Set(sc.args[0], v)

	var b bytes.Buffer
	b.WriteString(OK)

	return b.Bytes()
}

type GetCommand struct {
	BaseCommand
}

func NewGetCommand(args []string, flags []*Flag) *GetCommand {
	return &GetCommand{
		BaseCommand{
			args,
			flags,
		},
	}
}

func (gc *GetCommand) Execute(ds data.DataStore) []byte {
	log.Println("getting...")

	if len(gc.args) != 1 {
		log.Fatal(ErrInvalidNumberOfArguments)
	}

	v, ok := ds.Get(gc.args[0])
	if !ok {
		var b bytes.Buffer
		b.WriteString(NULL_BULK_STRING)
		return b.Bytes()
	}

	res := v.(*data.RedisValue)
	return writeBulkString(res.Value())
}
