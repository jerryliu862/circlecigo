package mailServer

import (
	"17live_wso_be/util"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func (c *Client) SendNoRegionNotification(ctx context.Context, mailList []string, campaignList []int) error {
	content := c.Content.NoRegion

	for _, address := range mailList {
		a := strings.Split(address, "@")
		if len(a) != 2 {
			log.Errorf("invalid email address: %s", address)
			return fmt.Errorf("invalid email address %s", address)
		}

		campaignIds := util.IntSliceToString(campaignList)

		plainTextContent := strings.Replace(content.PlainText, "?", campaignIds, 1)
		htmlContent := strings.Replace(content.HtmlText, "?", campaignIds, 1)

		if err := c.send(ctx, address, a[0], content.Subject, plainTextContent, htmlContent); err != nil {
			return err
		}

		log.Infof("no region notification sent: %s", address)
	}

	return nil
}

func (c *Client) SendSyncDataFinishNotification(ctx context.Context, emailAddress string) error {
	content := c.Content.SyncDataFinish

	a := strings.Split(emailAddress, "@")
	if len(a) != 2 {
		log.Errorf("invalid email address: %s", emailAddress)
		return fmt.Errorf("invalid email address %s", emailAddress)
	}

	if err := c.send(ctx, emailAddress, a[0], content.Subject, content.PlainText, content.HtmlText); err != nil {
		return err
	}

	log.Infof("sync data finish notification sent: %s", emailAddress)

	return nil
}

func (c *Client) send(ctx context.Context, receiver, receiverName, subject, plainTextContent, htmlContent string) error {
	from := mail.NewEmail(c.SenderName, c.Sender)
	to := mail.NewEmail(receiverName, receiver)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(c.ApiKey)

	if resp, err := client.Send(message); err != nil {
		log.Errorf("fail to send email via sendgrid: reveiver %s, %s", receiver, err.Error())
		return err
	} else if !util.ContainInt([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		log.Errorf("got unexpected status code when send email to %s: %d, %s", receiver, resp.StatusCode, resp.Body)
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return nil
}
