package common

import (
	"encoding/json"
	"github.com/JeffreyRichter/enum/enum"
	"reflect"
	"strings"
	"time"
)

var eOutputMessageType = outputMessageType(0)

// outputMessageType defines the nature of the output, ex: progress report, job summary, or error
type outputMessageType uint8

func (outputMessageType) Init() outputMessageType     { return outputMessageType(0) } // simple print, allowed to float up
func (outputMessageType) Info() outputMessageType     { return outputMessageType(1) } // simple print, allowed to float up
func (outputMessageType) Progress() outputMessageType { return outputMessageType(2) } // should be printed on the same line over and over again, not allowed to float up
func (outputMessageType) Exit() outputMessageType     { return outputMessageType(3) } // exit after printing
func (outputMessageType) Error() outputMessageType    { return outputMessageType(4) } // indicate fatal error, exit right after
func (outputMessageType) Prompt() outputMessageType   { return outputMessageType(5) } // ask the user a question after erasing the progress

func (o outputMessageType) String() string {
	return enum.StringInt(o, reflect.TypeOf(o))
}

// defines the output and how it should be handled
type outputMessage struct {
	msgContent   string
	msgType      outputMessageType
	exitCode     ExitCode      // only for when the application is meant to exit after printing (i.e. Error or Final)
	inputChannel chan<- string // support getting a response from the user
}

// used for output types that are not simple strings, such as progress and init
// a given format(text,json) is passed in, and the appropriate string is returned
type OutputBuilder func(OutputFormat) string

// -------------------------------------- JSON templates -------------------------------------- //
// used to help formatting of JSON outputs

func GetJsonStringFromTemplate(template interface{}) string {
	jsonOutput, err := json.Marshal(template)
	PanicIfErr(err)

	return string(jsonOutput)
}

// defines the general output template when the format is set to json
type jsonOutputTemplate struct {
	TimeStamp      time.Time
	MessageType    string
	MessageContent string // a simple string for INFO and ERROR, a serialized JSON for INIT, PROGRESS, EXIT
}

func newJsonOutputTemplate(messageType outputMessageType, messageContent string) *jsonOutputTemplate {
	return &jsonOutputTemplate{TimeStamp: time.Now(), MessageType: messageType.String(), MessageContent: messageContent}
}

type InitMsgJsonTemplate struct {
	LogFileLocation string
	JobID           string
}

func GetStandardInitOutputBuilder(jobID string, logFileLocation string) OutputBuilder {
	return func(format OutputFormat) string {
		if format == EOutputFormat.Json() {
			return GetJsonStringFromTemplate(InitMsgJsonTemplate{
				JobID:           jobID,
				LogFileLocation: logFileLocation,
			})
		}

		var sb strings.Builder
		sb.WriteString("\nJob " + jobID + " has started\n")
		sb.WriteString("Log file is located at: " + logFileLocation)
		sb.WriteString("\n")
		return sb.String()
	}
}
