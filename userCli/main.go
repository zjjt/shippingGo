package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/cli"
	microclient "github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/cmd"
	pb "github.com/zjjt/shippingGo/userService/proto/user"
)

func main() {
	cmd.Init()
	client := pb.NewUserServiceClient("shippingGo.service.user", microclient.DefaultClient)
	//utilisation de la ligne de commande pour passer des parametres via l'outil de go-micro
	service := micro.NewService(
		micro.Flags(
			cli.StringFlag{
				Name:  "name",
				Usage: "Your full name",
			},
			cli.StringFlag{
				Name:  "email",
				Usage: "Your email",
			},
			cli.StringFlag{
				Name:  "password",
				Usage: "Your password",
			},
			cli.StringFlag{
				Name:  "company",
				Usage: "Your company",
			},
		),
	)
	//start the service
	service.Init(
		micro.Action(func(c *cli.Context) {
			name := c.String("name")
			email := c.String("email")
			password := c.String("email")
			company := c.String("company")

			r, err := client
		}),
	)

}
