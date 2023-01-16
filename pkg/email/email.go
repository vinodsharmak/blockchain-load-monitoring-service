package email

import (
	"errors"
	"strings"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/constants"
	"github.com/antigloss/go/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var log *logger.Logger
var service *ses.SES

// Creating session for using service
func Config() error {
	log = config.Logger

	session, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		return err
	}

	service = ses.New(session)

	return nil
}

func SendEmail(message string) error {
	log := config.Logger
	log.Info("Email service start")

	sender := config.SenderEmail
	receivers := strings.Split(config.ReceiverEmails, ",")
	subject := "Blockchain Monitoring Service Alert !"
	charSet := constants.EmailCharSet

	//validate data
	if message == "" {
		return errors.New("message cannot not be blank")
	}

	// Preparing input for SendEmail Function
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: aws.StringSlice(receivers),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(message),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}

	// Sending Email
	result, err := service.SendEmail(input)

	if err != nil {
		if err1, ok := err.(awserr.Error); ok {
			switch err1.Code() {
			case ses.ErrCodeMessageRejected:
				log.Infof(ses.ErrCodeMessageRejected, err1.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				log.Infof(ses.ErrCodeMailFromDomainNotVerifiedException, err1.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				log.Infof(ses.ErrCodeConfigurationSetDoesNotExistException, err1.Error())
			default:
				log.Infof(err1.Error())
			}
			return err1
		} else {
			return err
		}
	}
	log.Info("Email sent with message ID : ", result.MessageId)
	log.Info("Email service end")
	return err
}
