package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Token string
)

type weapon struct {
	damage  int
	defense int
}

type field struct {
	monsterKind string
	fosition_x  int
	foistion_y  int
	attribute   int // 0 : fight 1 : shop
}

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	discordConn, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating discord session,", err)
		return
	}
	// login
	discordConn.AddHandler(messagePrint)
	err = discordConn.Open()
	if err != nil {
		fmt.Println("error creating discord session,", err)
		return
	}
	fmt.Println("Bot is now running. Press CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discordConn.Close()
}

func messagePrint(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	// user start
	username := m.Author.Username
	userid := m.Author.ID
	// usr end

	// map start
	forest := field{monsterKind: "slime", fosition_x: 0, foistion_y: 0, attribute: 0}
	// map end

	// weapon start
	stick := weapon{damage: 2, defense: 1}       // 나뭇가지
	bambooSpear := weapon{damage: 5, defense: 0} // 죽창
	stone := weapon{damage: 2, defense: 8}       // 짱돌
	cane := weapon{damage: 3, defense: 5}        // 지팡이
	rustySword := weapon{damage: 7, defense: 2}  // 녹슨 검
	// weapon end

	if m.Content == "!"+"help" {
		s.ChannelMessageSend(m.ChannelID, "이론Bot의 명령어 모음\n\n\n!help : 도움말을 출력합니다.\n\n\n그냥 만들어둔 기능\n\n!youtubelink : 그냥 개발자가 추천하는 유튜브 영상 링크 출력\n\n!img : 그냥 개발자가 아무거나 올린 이미지 출력\n\n\n게임 명령어\n\n!login : 로그인을 실행합니다. \n만약 처음 로그인을 한다면 자동 로그인을 실시합니다.\n\n!logout : 로그아웃 합니다.\n\n!weaponlist : 현재 게임상 존재하는 무기를 출력합니다.\n\n!status : 본인 계정의 레벨과 경험치를 출력합니다.\n\n\n")
	}

	// just start
	if m.Content == "!"+"youtubelink" {
		s.ChannelMessageSend(m.ChannelID, "https://www.youtube.com/watch?v=earHqoVE4HY")
	}
	if m.Content == "!"+"img" {
		s.ChannelMessageSend(m.ChannelID, "https://cdn.discordapp.com/attachments/842760181447131187/843145552324853770/7c4038d8061f97ea.jpg")
	}
	// just end

	// game start
	if m.Content == "!"+"weaponlist" {
		var conncheck int
		conn, err := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "데이터 베이스 연결에 실패했습니다.")
			return
		}
		err = conn.QueryRow("select conncheck from user where userid = ?", userid).Scan(&conncheck)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error")
			return
		}
		// "\n(weaponName)\n\n공격력 = "+transTypeIntToString((weaponName).damage)+"\n방어력 = "+transTypeIntToString((weaponName).defense)+"\n--------------------"
		if conncheck == 1 {
			s.ChannelMessageSend(m.ChannelID, "무기 종류\n\n--------------------"+"\n나뭇가지\n\n공격력 = "+transTypeIntToString(stick.damage)+"\n방어력 = "+transTypeIntToString(stick.defense)+"\n--------------------"+"\n죽창\n\n공격력 = "+transTypeIntToString(bambooSpear.damage)+"\n방어력 = "+transTypeIntToString(bambooSpear.defense)+"\n--------------------"+"\n짱돌\n\n공격력 = "+transTypeIntToString(stone.damage)+"\n방어력 = "+transTypeIntToString(stone.defense)+"\n--------------------"+"\n지팡이\n\n공격력 = "+transTypeIntToString(cane.damage)+"\n방어력 = "+transTypeIntToString(cane.defense)+"\n--------------------"+"\n녹슨 검\n\n공격력 = "+transTypeIntToString(rustySword.damage)+"\n방어력 = "+transTypeIntToString(rustySword.defense)+"\n--------------------")
		} else {
			s.ChannelMessageSend(m.ChannelID, "로그인하지 않았습니다.")
		}
		conn.Close()
	}
	if m.Content == "!"+"myfield" {
		var conncheck int
		var user_field string
		conn, err := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "데이터 베이스 연결에 실패했습니다.")
			return
		}
		err = conn.QueryRow("select conncheck from user where userid = ?", userid).Scan(&conncheck)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error")
			return
		}
		if conncheck == 1 {
			var userx int
			var usery int
			err = conn.QueryRow("select fieldname from userpos where userid = ?", userid).Scan(&user_field)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error")
				return
			}
			err = conn.QueryRow("select userx from userpos where userid = ?", userid).Scan(&userx)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error")
				return
			}
			err = conn.QueryRow("select usery from userpos where userid = ?", userid).Scan(&usery)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error")
				return
			}
			s.ChannelMessageSend(m.ChannelID, "나의 위치\n\n--------------------\n")
			if userx == 0 && usery == 0 {
				s.ChannelMessageSend(m.ChannelID, user_field+"\n\n나오는 몬스터 = "+forest.monsterKind+"\n\n위치\n(! 가 유저의 위치입니다.)\n"+mapString(forest.fosition_x, forest.attribute)+"\n--------------------")
			} else {
				var userx int
				var usery int
				err = conn.QueryRow("select userx from userpos where userid = ?", userid).Scan(&userx)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "error")
					return
				}
				err = conn.QueryRow("select usery from userpos where userid = ?", userid).Scan(&usery)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "error")
					return
				}
				s.ChannelMessageSend(m.ChannelID, "주요한 맵이 아니므로 정보가 없습니다.\n"+mapString(userx, usery)+"\n--------------------")
			}
		}
		conn.Close()
	}
	if m.Content == "!"+"right" {
		var conncheck int
		conn, err := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "데이터 베이스 연결에 실패했습니다.")
			return
		}
		err = conn.QueryRow("select conncheck from user where userid = ?", userid).Scan(&conncheck)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error")
			return
		}
		if conncheck == 1 {
			var userx int
			err = conn.QueryRow("select userx from userpos where userid = ?", userid).Scan(&userx)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error")
				return
			}
			if userx < 9 {
				changex := userx + 1
				upd, _ := conn.Exec("update userpos set userx = ? where userid = ?", changex, userid)
				upd_check, _ := upd.RowsAffected()
				var change_userx int
				var usery int
				err = conn.QueryRow("select userx from userpos where userid = ?", userid).Scan(&change_userx)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "error")
					return
				}
				err = conn.QueryRow("select usery from userpos where userid = ?", userid).Scan(&usery)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "error")
					return
				}
				if upd_check == 1 {
					s.ChannelMessageSend(m.ChannelID, username+"의 위치\n(!가 현재 위치입니다.)\n--------------------\n"+mapString(change_userx, usery)+"\n--------------------")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "이동 불가능 지역입니다.")
			}
		}
		conn.Close()
	}
	if m.Content == "!"+"left" {
		var conncheck int
		conn, err := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "데이터 베이스 연결에 실패했습니다.")
			return
		}
		err = conn.QueryRow("select conncheck from user where userid = ?", userid).Scan(&conncheck)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error")
			return
		}
		if conncheck == 1 {
			var userx int
			err = conn.QueryRow("select userx from userpos where userid = ?", userid).Scan(&userx)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error")
				return
			}
			if userx > 0 {
				changex := userx - 1
				upd, _ := conn.Exec("update userpos set userx = ? where userid = ?", changex, userid)
				upd_check, _ := upd.RowsAffected()
				var change_userx int
				var usery int
				err = conn.QueryRow("select userx from userpos where userid = ?", userid).Scan(&change_userx)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "error")
					return
				}
				err = conn.QueryRow("select usery from userpos where userid = ?", userid).Scan(&usery)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "error")
					return
				}
				if upd_check == 1 {
					s.ChannelMessageSend(m.ChannelID, username+"의 위치\n(!가 현재 위치입니다.)\n--------------------\n"+mapString(change_userx, usery)+"\n--------------------")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "이동 불가능 지역입니다.")
			}
		}
		conn.Close()
	}
	if m.Content == "!"+"up" {
		var conncheck int
		conn, err := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "데이터 베이스 연결에 실패했습니다.")
			return
		}
		err = conn.QueryRow("select conncheck from user where userid = ?", userid).Scan(&conncheck)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error")
			return
		}
		if conncheck == 1 {
			var usery int
			err = conn.QueryRow("select usery from userpos where userid = ?", userid).Scan(&usery)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error")
				return
			}
			if usery > 0 {
				changey := usery - 1
				upd, _ := conn.Exec("update userpos set usery = ? where userid = ?", changey, userid)
				upd_check, _ := upd.RowsAffected()
				var userx int
				var change_usery int
				err = conn.QueryRow("select userx from userpos where userid = ?", userid).Scan(&userx)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "error")
					return
				}
				err = conn.QueryRow("select usery from userpos where userid = ?", userid).Scan(&change_usery)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "error")
					return
				}
				if upd_check == 1 {
					s.ChannelMessageSend(m.ChannelID, username+"의 위치\n(!가 현재 위치입니다.)\n--------------------\n"+mapString(userx, change_usery)+"\n--------------------")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "이동 불가능 지역입니다.")
			}
		}
		conn.Close()
	}
	if m.Content == "!"+"down" {
		var conncheck int
		conn, err := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "데이터 베이스 연결에 실패했습니다.")
			return
		}
		err = conn.QueryRow("select conncheck from user where userid = ?", userid).Scan(&conncheck)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error")
			return
		}
		if conncheck == 1 {
			var usery int
			err = conn.QueryRow("select usery from userpos where userid = ?", userid).Scan(&usery)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error")
				return
			}
			if usery < 9 {
				changey := usery + 1
				upd, _ := conn.Exec("update userpos set usery = ? where userid = ?", changey, userid)
				upd_check, _ := upd.RowsAffected()
				var userx int
				var change_usery int
				err = conn.QueryRow("select userx from userpos where userid = ?", userid).Scan(&userx)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "error")
					return
				}
				err = conn.QueryRow("select usery from userpos where userid = ?", userid).Scan(&change_usery)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "error")
					return
				}
				if upd_check == 1 {
					s.ChannelMessageSend(m.ChannelID, username+"의 위치\n(!가 현재 위치입니다.)\n--------------------\n"+mapString(userx, change_usery)+"\n--------------------")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "이동 불가능 지역입니다.")
			}
		}
		conn.Close()
	}
	if m.Content == "!"+"login" {
		var temp_userid string
		conn, err := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "데이터 베이스 연결에 실패했습니다.")
			return
		}
		err = conn.QueryRow("select userid from user where userid = ?", userid).Scan(&temp_userid)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "로그인 실패...")
			time.Sleep(1000)
			s.ChannelMessageSend(m.ChannelID, "자동 로그인을 실행합니다.")
			ins_1, err := conn.Exec("insert into user(userid) values(?)", userid)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "자동 로그인에 실패했습니다.")
				return
			}
			ins_2, err := conn.Exec("insert into status(userid) values(?)", userid)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "자동 로그인에 실패했습니다.")
				return
			}
			ins_3, err := conn.Exec("insert into userpos(userid) values(?)", userid)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "자동 로그인에 실패했습니다.")
			}
			ins_check_1, _ := ins_1.RowsAffected()
			ins_check_2, _ := ins_2.RowsAffected()
			ins_check_3, _ := ins_3.RowsAffected()
			if ins_check_1 == 1 && ins_check_2 == 1 && ins_check_3 == 1 {
				upd, _ := conn.Exec("update user set conncheck = 1 where userid = ?", userid)
				upd_check, _ := upd.RowsAffected()
				if upd_check == 1 {
					s.ChannelMessageSend(m.ChannelID, "자동 로그인에 성공하였습니다.")
				} else {
					s.ChannelMessageSend(m.ChannelID, "이미 로그인 되었습니다.")
				}
				upd_check = 0
			}
		} else {
			upd, _ := conn.Exec("update user set conncheck = 1 where userid = ?", userid)
			upd_check, _ := upd.RowsAffected()
			if upd_check == 1 {
				s.ChannelMessageSend(m.ChannelID, "로그인에 성공하였습니다.")
			} else {
				s.ChannelMessageSend(m.ChannelID, "이미 로그인 되었습니다.")
			}
		}
		conn.Close()
	}
	if m.Content == "!"+"logout" {
		var conncheck int
		conn, err := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "데이터 베이스 연결에 실패했습니다.")
			return
		}
		err = conn.QueryRow("select conncheck from user where userid = ?", userid).Scan(&conncheck)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error")
			return
		}
		if conncheck == 0 {
			s.ChannelMessageSend(m.ChannelID, "아직 로그인하지 않았습니다.")
			return
		} else if conncheck == 1 {
			upd, _ := conn.Exec("update user set conncheck = 0 where userid = ?", userid)
			upd_check, _ := upd.RowsAffected()
			if upd_check == 1 {
				s.ChannelMessageSend(m.ChannelID, "로그아웃에 성공하였습니다.")
			}
			upd_check = 0
		}
		conn.Close()
	}
	if m.Content == "!"+"status" {
		var conncheck int
		conn, err := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "데이터 베이스 연결에 실패했습니다.")
			return
		}
		err = conn.QueryRow("select conncheck from user where userid = ?", userid).Scan(&conncheck)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error")
			return
		}
		if conncheck == 1 {
			var level, exp string
			err = conn.QueryRow("select level from status where userid = ?", userid).Scan(&level)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "스테이터스를 불러들일 수 없습니다.")
				return
			}
			err = conn.QueryRow("select exp from status where userid = ?", userid).Scan(&exp)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "스테이터스를 불러들일 수 없습니다.")
				return
			}
			s.ChannelMessageSend(m.ChannelID, username+" 님의 레벨 : "+level+" LV")
			s.ChannelMessageSend(m.ChannelID, username+" 님의 경험치 : "+exp+" EXP")
		} else {
			s.ChannelMessageSend(m.ChannelID, "로그인 하지 않았습니다.")
		}
	}
	// game end
}

func transTypeIntToString(n int) string {
	str := ""
	if n < 10 {
		str = str + string(n+48)
	}
	return str
}

func mapString(x int, y int) string {
	var strArr [10][61]byte
	str := ""
	for i := 0; i < 10; i++ {
		for j := 0; j < 60; j++ {
			if j%6 == 0 {
				strArr[i][j] = '.'
			} else {
				strArr[i][j] = ' '
			}
			if i == y && j == (x*6) {
				strArr[i][j] = '!'
			}
			str = str + string(strArr[i][j])
		}
		strArr[i][60] = '\n'
		str = str + string(strArr[i][60])
	}
	return str
}
