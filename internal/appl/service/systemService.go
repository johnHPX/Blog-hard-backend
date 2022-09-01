package service

import (
	"fmt"

	"github.com/johnHPX/blog-hard-backend/internal/infra/repository"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/configsAPI"
	mail "github.com/xhit/go-simple-mail/v2"
)

type systemServiceInterface interface {
	SendEmail(template, emailToDestiny, messageTitle string) error
	SendEmailComment(CommentIdOrResponseCommentId string, commentOrResponseComment bool) error
}

type systemServiceImpl struct{}

func (S *systemServiceImpl) SendEmail(template, emailToDestiny, messageTitle string) error {
	// gets configs contact
	configService := configsAPI.NewConfigs()
	contactConfig, err := configService.ContactConfig()
	if err != nil {
		return err
	}

	// create a client smtp
	server := mail.NewSMTPClient()
	server.Host = "smtp.gmail.com"
	server.Port = 587
	server.Username = contactConfig.Email
	server.Password = contactConfig.Secret
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	// Create email
	emailSend := mail.NewMSG()
	emailSend.SetFrom(fmt.Sprintf("From Me <%s>", contactConfig.Email))
	emailSend.AddTo(emailToDestiny)
	emailSend.SetSubject(messageTitle)
	emailSend.SetBody(mail.TextHTML, template)

	// Send email
	err = emailSend.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

func (s *systemServiceImpl) SendEmailComment(CommentIdOrResponseCommentId string, commentOrResponseComment bool) error {

	// if true, is a comment of post
	if commentOrResponseComment {
		repComment := repository.NewCommentRepository()
		commentEntity, err := repComment.Find(CommentIdOrResponseCommentId)
		if err != nil {
			return err
		}

		repPost := repository.NewPostRepository()
		postEntity, err := repPost.Find(commentEntity.PostID)
		if err != nil {
			return err
		}

		repUser := repository.NewUserRepository()
		userEntity, err := repUser.Find(commentEntity.UserID)
		if err != nil {
			return err
		}

		template := fmt.Sprintf(`

		<!DOCTYPE html>
		<html lang="pt-br">
			<head>
    			<meta charset="UTF-8">
    			<title>Atualizações do site</title>
			</head>
			<body>
				<h1>Atualizações do Blog</h1>

				<p>Olá Adminstrador! Aqui está as novidades do seu blog!</p>

				<h3>Um Usuário comentou em um dos Seus Posts!</h3>
				<p>O usuario %s comentou no Post %s</p>
			</body>
		</html>

	`, userEntity.Nick, postEntity.Title)

		admUserEntities, err := repUser.ListUserAdm("adm", 0, 10, 1)
		if err != nil {
			return err
		}

		for _, v := range admUserEntities {
			err = s.SendEmail(template, v.Email, "Atualizações do Blog")
			if err != nil {
				return err
			}
		}
		return nil

	}

	// if false, is a response of comment of post
	repReponseComment := repository.NewResponseCommmentRepository()
	responsecommentEntity, err := repReponseComment.Find(CommentIdOrResponseCommentId)
	if err != nil {
		return err
	}

	repComment := repository.NewCommentRepository()
	commentEntity, err := repComment.Find(responsecommentEntity.CommentID)
	if err != nil {
		return err
	}

	repPost := repository.NewPostRepository()
	postEntity, err := repPost.Find(commentEntity.PostID)
	if err != nil {
		return err
	}

	repUser := repository.NewUserRepository()
	userEntityResponseComment, err := repUser.Find(responsecommentEntity.UserID)
	if err != nil {
		return err
	}

	userEntityComment, err := repUser.Find(commentEntity.UserID)
	if err != nil {
		return err
	}

	template1 := fmt.Sprintf(`

	<!DOCTYPE html>
	<html lang="pt-br">
		<head>
			<meta charset="UTF-8">
			<title>Atualizações do site</title>
		</head>
		<body>
			<h1>Atualizações do Blog</h1>

			<p>Olá Adminstrador! Aqui está as novidades do seu blog!</p>

			<h3>Um Usuário respondeu ao um comentario de um dos Seus Posts!</h3>
			<p>O usuario %s respondeu ao comentario de %s no Post %s</p>
		</body>
	</html>

`, userEntityResponseComment.Nick, userEntityComment.Nick, postEntity.Title)

	admUserEntities, err := repUser.ListUserAdm("adm", 0, 10, 1)
	if err != nil {
		return err
	}

	for _, v := range admUserEntities {
		err = s.SendEmail(template1, v.Email, "Atualizações do Blog")
		if err != nil {
			return err
		}
	}

	return nil

}

func NewSystemService() systemServiceInterface {
	return &systemServiceImpl{}
}
