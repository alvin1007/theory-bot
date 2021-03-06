package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
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

type monster struct {
	monstername  string
	hp           int
	damage       int
	shield       int
	money        int
	exp          int
	weaponname   []string
	weaponrandom []int
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

	//server start
	serverid := m.GuildID
	//server end

	// map start
	forest := field{monsterKind: "slime", fosition_x: 0, foistion_y: 0, attribute: 0}
	// map end

	// weapon start
	// stage 1 start
	var weaponname1 = []string{"stick", "bambooSpear", "stone", "cane", "rustySword"}
	var weaponrandom1 = []int{50, 30, 10, 7, 3}
	stick := weapon{damage: 2, defense: 1}       // 나뭇가지
	bambooSpear := weapon{damage: 5, defense: 0} // 죽창
	stone := weapon{damage: 2, defense: 8}       // 짱돌
	cane := weapon{damage: 3, defense: 5}        // 지팡이
	rustySword := weapon{damage: 7, defense: 2}  // 녹슨 검
	//stage 1 end
	// weapon end

	// monster start
	slime := monster{monstername: "슬라임", hp: 50, damage: 5, shield: 2, money: 1, exp: 5, weaponname: weaponname1, weaponrandom: weaponrandom1}
	// monster end

	slice := strings.Split(m.Content, " ")

	if m.Content == cognition(serverid)+"help" {
		s.ChannelMessageSend(m.ChannelID, "이론Bot의 명령어 모음\n\n\n"+cognition(serverid)+"help : 도움말을 출력합니다.\n\n\n그냥 만들어둔 기능\n\n"+cognition(serverid)+"youtubelink : 그냥 개발자가 추천하는 유튜브 영상 링크 출력\n\n"+cognition(serverid)+"img : 그냥 개발자가 아무거나 올린 이미지 출력\n\n\n게임 명령어\n\n"+cognition(serverid)+"login : 로그인을 실행합니다. \n만약 처음 로그인을 한다면 자동 로그인을 실시합니다.\n\n"+cognition(serverid)+"logout : 로그아웃 합니다.\n\n"+cognition(serverid)+"weaponlist : 현재 게임상 존재하는 무기를 출력합니다.\n\n"+cognition(serverid)+"status : 본인 계정의 레벨과 경험치를 출력합니다.\n\n\n")
	}

	// just start
	if m.Content == cognition(serverid)+"youtubelink" {
		s.ChannelMessageSend(m.ChannelID, "https://youtu.be/jz_Jozym5DM")
	}
	if m.Content == cognition(serverid)+"img" {
		s.ChannelMessageSend(m.ChannelID, "https://cdn.discordapp.com/attachments/843111989735063585/844564052205633626/150619_5fb4b9dbde844.png")
	}
	if slice[0] == cognition(serverid)+"changecog" {
		conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
		upd, err := conn.Exec("update cognition set cognition = ? where serverid = ?", slice[1], serverid)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "변경 실패했습니다.")
		}
		upd_check, _ := upd.RowsAffected()
		if upd_check == 1 {
			s.ChannelMessageSend(m.ChannelID, "정상적으로 변경되었습니다.")
		} else {
			s.ChannelMessageSend(m.ChannelID, "이미 설정된 인식어입니다.")
		}
		conn.Close()
	}
	// just end

	// anime-brain start
	if slice[0] == cognition(serverid)+"bf" {
		s.ChannelMessageSend(m.ChannelID, bf(slice[1]))
	}
	if slice[0] == cognition(serverid)+"genshinbf" {
		s.ChannelMessageSend(m.ChannelID, genshin_bf(slice))
	}
	if slice[0] == cognition(serverid)+"bftogenshin" {
		s.ChannelMessageSend(m.ChannelID, bfToGenshinbf(slice[1]))
	}
	// anime-brain end

	// game start
	if m.Content == cognition(serverid)+"weaponlist" {
		// "\n(weaponName)\n\n공격력 = "+transTypeIntToString((weaponName).damage)+"\n방어력 = "+transTypeIntToString((weaponName).defense)+"\n--------------------"
		if loginCheck(userid) {
			s.ChannelMessageSend(m.ChannelID, "무기 종류\n\n--------------------"+"\n나뭇가지\n\n공격력 = "+transTypeIntToString(stick.damage)+"\n방어력 = "+transTypeIntToString(stick.defense)+"\n--------------------"+"\n죽창\n\n공격력 = "+transTypeIntToString(bambooSpear.damage)+"\n방어력 = "+transTypeIntToString(bambooSpear.defense)+"\n--------------------"+"\n짱돌\n\n공격력 = "+transTypeIntToString(stone.damage)+"\n방어력 = "+transTypeIntToString(stone.defense)+"\n--------------------"+"\n지팡이\n\n공격력 = "+transTypeIntToString(cane.damage)+"\n방어력 = "+transTypeIntToString(cane.defense)+"\n--------------------"+"\n녹슨 검\n\n공격력 = "+transTypeIntToString(rustySword.damage)+"\n방어력 = "+transTypeIntToString(rustySword.defense)+"\n--------------------")
		} else {
			s.ChannelMessageSend(m.ChannelID, "로그인하지 않았습니다.")
		}
	}
	if m.Content == cognition(serverid)+"myfield" {
		if loginCheck(userid) {
			var userx int
			var usery int
			conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
			err := conn.QueryRow("select userx from userpos where userid = ?", userid).Scan(&userx)
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
			if userx == forest.fosition_x && usery == forest.foistion_y {
				s.ChannelMessageSend(m.ChannelID, "포레스트\n\n나오는 몬스터 = "+slime.monstername+"\n\n위치\n(! 가 유저의 위치입니다.)\n"+mapString(forest.fosition_x, forest.foistion_y)+"\n--------------------")
			} else {
				var user_x int
				var user_y int
				err = conn.QueryRow("select userx from userpos where userid = ?", userid).Scan(&user_x)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "error")
					return
				}
				err = conn.QueryRow("select usery from userpos where userid = ?", userid).Scan(&user_y)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "error")
					return
				}
				s.ChannelMessageSend(m.ChannelID, "주요한 맵이 아니므로 정보가 없습니다.\n"+mapString(user_x, user_y)+"\n--------------------")
			}
			conn.Close()
		} else {
			s.ChannelMessageSend(m.ChannelID, "로그인을 먼저 해야합니다.")
		}
	}
	if m.Content == cognition(serverid)+"right" {
		if loginCheck(userid) {
			var userx int
			conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
			err := conn.QueryRow("select userx from userpos where userid = ?", userid).Scan(&userx)
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
				s.ChannelMessageSend(m.ChannelID, "이동 불가능 지역입니다.\n\n--------------------\n"+"현재 "+username+"의 위치입니다.\n\n"+mapString(userx, usery)+"\n--------------------")
			}
			conn.Close()
		} else {
			s.ChannelMessageSend(m.ChannelID, "로그인을 먼저 해야합니다.")
		}
	}
	if m.Content == cognition(serverid)+"left" {
		if loginCheck(userid) {
			var userx int
			conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
			err := conn.QueryRow("select userx from userpos where userid = ?", userid).Scan(&userx)
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
				s.ChannelMessageSend(m.ChannelID, "이동 불가능 지역입니다.\n\n--------------------\n"+"현재 "+username+"의 위치입니다.\n\n"+mapString(userx, usery)+"\n--------------------")
			}
			conn.Close()
		} else {
			s.ChannelMessageSend(m.ChannelID, "로그인을 먼저 해야합니다.")
		}
	}
	if m.Content == cognition(serverid)+"up" {
		if loginCheck(userid) {
			var usery int
			conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
			err := conn.QueryRow("select usery from userpos where userid = ?", userid).Scan(&usery)
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
				s.ChannelMessageSend(m.ChannelID, "이동 불가능 지역입니다.\n\n--------------------\n"+"현재 "+username+"의 위치입니다.\n\n"+mapString(userx, usery)+"\n--------------------")
			}
			conn.Close()
		} else {
			s.ChannelMessageSend(m.ChannelID, "로그인을 먼저 해야합니다.")
		}
	}
	if m.Content == cognition(serverid)+"down" {
		if loginCheck(userid) {
			var usery int
			conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
			err := conn.QueryRow("select usery from userpos where userid = ?", userid).Scan(&usery)
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
				s.ChannelMessageSend(m.ChannelID, "이동 불가능 지역입니다.\n\n--------------------\n"+"현재 "+username+"의 위치입니다.\n\n"+mapString(userx, usery)+"\n--------------------")
			}
			conn.Close()
		} else {
			s.ChannelMessageSend(m.ChannelID, "로그인 하지 않았습니다.")
		}
	}
	if m.Content == cognition(serverid)+"login" {
		var temp_userid string
		conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
		err := conn.QueryRow("select userid from user where userid = ?", userid).Scan(&temp_userid)
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
			ins_4, err := conn.Exec("insert into fight(userid) values(?)", userid)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "자동 로그인에 실패했습니다.")
			}
			ins_5, err := conn.Exec("insert into usermoney(userid) values(?)", userid)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "자동 로그인에 실패했습니다.")
			}
			ins_check_1, _ := ins_1.RowsAffected()
			ins_check_2, _ := ins_2.RowsAffected()
			ins_check_3, _ := ins_3.RowsAffected()
			ins_check_4, _ := ins_4.RowsAffected()
			ins_check_5, _ := ins_5.RowsAffected()
			if ins_check_1 == 1 && ins_check_2 == 1 && ins_check_3 == 1 && ins_check_4 == 1 && ins_check_5 == 1 {
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
	if m.Content == cognition(serverid)+"logout" {
		if loginCheck(userid) {
			conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
			upd, _ := conn.Exec("update user set conncheck = 0 where userid = ?", userid)
			upd_check, _ := upd.RowsAffected()
			if upd_check == 1 {
				s.ChannelMessageSend(m.ChannelID, "로그아웃에 성공하였습니다.")
			}
			upd_check = 0
			conn.Close()
		} else {
			s.ChannelMessageSend(m.ChannelID, "로그인 하지 않았습니다.")
		}
	}
	if m.Content == cognition(serverid)+"status" {
		if loginCheck(userid) {
			conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
			var level, exp string
			var money string
			err := conn.QueryRow("select level from status where userid = ?", userid).Scan(&level)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "스테이터스를 불러들일 수 없습니다.")
				return
			}
			err = conn.QueryRow("select exp from status where userid = ?", userid).Scan(&exp)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "스테이터스를 불러들일 수 없습니다.")
				return
			}
			err = conn.QueryRow("select money from usermoney where userid = ?", userid).Scan(&money)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "스테이터스를 불러들일 수 없습니다.")
				return
			}
			s.ChannelMessageSend(m.ChannelID, username+" 님의 레벨 : "+level+" LV")
			s.ChannelMessageSend(m.ChannelID, username+" 님의 경험치 : "+exp+" EXP")
			s.ChannelMessageSend(m.ChannelID, username+" 님의 골드 : "+money+" 골드")
			s.ChannelMessageSend(m.ChannelID, username+" 님의 공격력 : "+transTypeIntToString(userDamage(userid))+" DAMAGE")
			s.ChannelMessageSend(m.ChannelID, username+" 님의 방어력 : "+transTypeIntToString(userDefense(userid))+" DEFENSE")
			conn.Close()
		} else {
			s.ChannelMessageSend(m.ChannelID, "로그인 하지 않았습니다.")
		}
	}
	if m.Content == cognition(serverid)+"myweapon" {
		if loginCheck(userid) {
			userWeapon := userWeaponName(userid)
			if userWeapon[0] == "" {
				s.ChannelMessageSend(m.ChannelID, "무기가 없습니다.")
			}
			for i := 0; i < len(userWeapon); i++ {

			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "로그인 하지 않았습니다.")
		}
	}
	if m.Content == cognition(serverid)+"goavt" {
		if loginCheck(userid) {
			var userx int
			var usery int
			var fieldname string
			conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
			err := conn.QueryRow("select userx from userpos where userid = ?", userid).Scan(&userx)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error")
				return
			}
			err = conn.QueryRow("select usery from userpos where userid = ?", userid).Scan(&usery)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error")
				return
			}
			fieldname = fieldName(userx, usery)
			if fieldname == "실패" {
				s.ChannelMessageSend(m.ChannelID, "모험을 떠날 수 있는 지역이 아닙니다.")
				return
			}
			s.ChannelMessageSend(m.ChannelID, fieldname+"로 모험을 떠나는 중입니다...")
			time.Sleep(5000)
			rand.Seed(time.Now().UnixNano())
			random := rand.Intn(100)
			if random < 60 {
				upd, _ := conn.Exec("update fight set fight = 1 where userid = ?", userid)
				upd_check, _ := upd.RowsAffected()
				if upd_check == 1 {
					s.ChannelMessageSend(m.ChannelID, "전투가 시작되었습니다!\n\n")
				}
				// 아직 미구현
				upd, _ = conn.Exec("update fight set fight = 0 where userid = ?", userid)
				upd_check, _ = upd.RowsAffected()
				if upd_check == 1 {
					s.ChannelMessageSend(m.ChannelID, "전투가 종료되었습니다!\n\n")
				}
			} else if random >= 60 && random < 90 {
				rand.Seed(time.Now().UnixNano())
				money_random := rand.Intn(50)
				var user_money_random int
				err := conn.QueryRow("select money from usermoney where userid = ?", userid).Scan(&user_money_random)
				if err != nil {
					log.Fatal(err)
				}
				upd, _ := conn.Exec("update usermoney set money = ? where userid = ?", user_money_random+money_random, userid)
				upd_check, _ := upd.RowsAffected()
				if upd_check == 1 {
					s.ChannelMessageSend(m.ChannelID, "당신은 모험 중 "+transTypeIntToString(money_random)+" 골드를 획득하셨습니다.")
				}
			} else if random >= 90 {
				userWeapon := weaponRandom(fieldname, userid)
				if userWeapon == "fail" {
					s.ChannelMessageSend(m.ChannelID, "무기 획득에 실패했습니다.")
					return
				} else if userWeapon == "w" {
					s.ChannelMessageSend(m.ChannelID, "이미 획득한 무기입니다.")
					return
				}
				s.ChannelMessageSend(m.ChannelID, userWeapon+"를 획득했습니다.")
			}
			conn.Close()
		} else {
			s.ChannelMessageSend(m.ChannelID, "로그인 하지 않았습니다.")
		}
	}
	// game end
}

func fieldName(x int, y int) string {
	// 여기는 맵을 추가할 때마다 추가해야함.
	if x == 0 && y == 0 {
		return "포레스트"
	} else {
		return "실패"
	}
}

func weaponRandom(userField string, userid string) string {
	//여기는 맵을 추가할 때마다 추가해야함
	conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
	// 후에 개발해야 될 것들
	if userField == "포레스트" {

		stick := weapon{damage: 2, defense: 1}       // 나뭇가지 70%
		bambooSpear := weapon{damage: 5, defense: 0} // 죽창 10%
		stone := weapon{damage: 2, defense: 8}       // 짱돌 5%
		cane := weapon{damage: 3, defense: 5}        // 지팡이 3%
		rustySword := weapon{damage: 7, defense: 2}  // 녹슨 검 2%

		rand.Seed(time.Now().UnixNano())
		random := rand.Intn(100)
		fmt.Println(random)
		if random < 70 {
			if weaponCheck("stick", userid) {
				_, err := conn.Exec("insert into userweapon(damage, defense, id, weaponname) values(?, ?, ?, ?)", stick.damage, stick.defense, userid, "stick")
				if err != nil {
					panic(err)
				}
				return "stick"
			}
		} else if random < 80 {
			if weaponCheck("bambooSpear", userid) {
				_, err := conn.Exec("insert into userweapon(damage, defense, id, weaponname) values(?, ?, ?, ?)", bambooSpear.damage, bambooSpear.defense, userid, "bambooSpear")
				if err != nil {
					panic(err)
				}
				return "bambooSpear"
			} else {
				return "w"
			}
		} else if random < 88 {
			if weaponCheck("stone", userid) {
				_, err := conn.Exec("insert into userweapon(damage, defense, id, weaponname) values(?, ?, ?, ?)", stone.damage, stone.defense, userid, "stone")
				if err != nil {
					panic(err)
				}
				return "stone"
			} else {
				return "w"
			}
		} else if random < 96 {
			if weaponCheck("cane", userid) {
				_, err := conn.Exec("insert into userweapon(damage, defense, id, weaponname) values(?, ?, ?, ?)", cane.damage, cane.defense, userid, "cane")
				if err != nil {
					panic(err)
				}
				return "cane"
			} else {
				return "w"
			}
		} else if random <= 100 {
			if weaponCheck("rustySword", userid) {
				_, err := conn.Exec("insert into userweapon(damage, defense, id, weaponname) values(?, ?, ?, ?)", rustySword.damage, rustySword.defense, userid, "rustySword")
				if err != nil {
					panic(err)
				}
				return "rustySword"
			} else {
				return "w"
			}
		}
	}
	conn.Close()
	return "fail"
}

func userWeaponName(userid string) []string {
	conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
	rows, _ := conn.Query("select weaponname from userweapon where id = ?", userid)
	i := 0
	var userWeapon [100]string
	var weapon string
	for rows.Next() {
		err := rows.Scan(&weapon)
		if err != nil {
			panic(err)
		}
		userWeapon[i] = weapon
	}
	conn.Close()
	return userWeapon[:]
}

func userDefense(userid string) int {
	conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
	var defense int
	var userWeaponDefense int
	var userlevel int
	err := conn.QueryRow("select level from status where userid = ?", userid).Scan(&userlevel)
	if err != nil {
		panic(err)
	}
	rows, err := conn.Query("select defense from userweapon where id = ?", userid)
	if err != nil {
		userWeaponDefense = 0
	}
	for rows.Next() {
		err := rows.Scan(&userWeaponDefense)
		if err != nil {
			panic(err)
		}
		defense += userWeaponDefense
	}
	defense += userlevel*2 + 5
	conn.Close()
	return defense
}

func userDamage(userid string) int {
	conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
	var damage int
	var userWeaponDamage int
	var userlevel int
	err := conn.QueryRow("select level from status where userid = ?", userid).Scan(&userlevel)
	if err != nil {
		panic(err)
	}
	rows, err := conn.Query("select damage from userweapon where id = ?", userid)
	if err != nil {
		userWeaponDamage = 0
	}
	for rows.Next() {
		err := rows.Scan(&userWeaponDamage)
		if err != nil {
			panic(err)
		}
		damage += userWeaponDamage
	}
	damage += userlevel*2 + 5
	conn.Close()
	return damage
}

func weaponCheck(weaponname string, userid string) bool {
	conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
	var weapon string
	rows, err := conn.Query("select weaponname from userweapon where id = ?", userid)
	if err != nil {
		conn.Close()
		return false
	}
	for rows.Next() {
		err := rows.Scan(&weapon)
		if err != nil {
			log.Fatal(err)
		}
		if weapon == weaponname {
			conn.Close()
			return false
		}
	}
	conn.Close()
	return true
}

func loginCheck(userid string) bool {
	var conncheck int
	conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
	err := conn.QueryRow("select conncheck from user where userid = ?", userid).Scan(&conncheck)
	if err == nil {
		if conncheck == 1 {
			conn.Close()
			return true
		} else {
			conn.Close()
			return false
		}
	}
	return false
}

func transTypeIntToString(n int) string {
	str := ""
	if n < 10 {
		str = str + string(rune(n+48))
	} else if n < 100 {
		n1 := n % 10
		n2 := n / 10
		str = string(rune(n2+48)) + string(rune(n1+48))
	} else if n < 1000 {
		n2 := n % 100
		n1 := n % 10
		n3 := n / 100
		str = string(rune(n3+48)) + string(rune(n2+48)) + string(rune(n1+48))
	}
	return str
}

func cognition(serverid string) string {
	var cognition string
	conn, _ := sql.Open("mysql", "root:alvin1007@tcp(localhost:3306)/game")
	err := conn.QueryRow("select cognition from cognition where serverid = ?", serverid).Scan(&cognition)
	if err != nil {
		_, err = conn.Exec("insert into cognition(serverid) values(?)", serverid)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return "err"
		}
	}
	conn.Close()
	return cognition
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

func bfToGenshinbf(code string) string {
	cnt := 0
	output := ""
	for {
		if cnt == len(code) {
			break
		}
		if string(code[cnt]) == "+" {
			output += "호두 "
		} else if string(code[cnt]) == "-" {
			output += "벤티 "
		} else if string(code[cnt]) == ">" {
			output += "클레 "
		} else if string(code[cnt]) == "<" {
			output += "유라 "
		} else if string(code[cnt]) == "[" {
			output += "감우 "
		} else if string(code[cnt]) == "]" {
			output += "눈나 "
		} else if string(code[cnt]) == "." {
			output += "헤으응 "
		}
		cnt++
	}
	return output
}

func bf(code string) string {
	output := ""
	cnt := 0
	ptr := 0
	var ascii [30000]int
	for {
		if cnt == len(code) {
			break
		}
		if string(code[cnt]) == "+" {
			ascii[ptr] += 1
		} else if string(code[cnt]) == "-" {
			ascii[ptr] -= 1
		} else if string(code[cnt]) == ">" {
			ptr += 1
		} else if string(code[cnt]) == "<" {
			ptr -= 1
		} else if string(code[cnt]) == "[" {
			if ascii[ptr] == 0 {
				cnt1 := 1
				for cnt1 != 0 {
					cnt++
					if string(code[cnt]) == "[" {
						cnt1++
					} else if string(code[cnt]) == "]" {
						cnt1--
					} else {
						continue
					}
				}
			}
		} else if string(code[cnt]) == "]" {
			if ascii[ptr] != 0 {
				cnt2 := 1
				for cnt2 != 0 {
					cnt--
					if string(code[cnt]) == "]" {
						cnt2++
					} else if string(code[cnt]) == "[" {
						cnt2--
					} else {
						continue
					}
				}
			}
		} else if string(code[cnt]) == "." {
			output += string(ascii[ptr])
		}
		cnt++
	}
	return output
}

func genshin_bf(code []string) string {
	output := ""
	cnt := 0
	ptr := 0
	var ascii [30000]int
	for {
		if cnt == len(code) {
			break
		}
		if code[cnt] == "호두" {
			ascii[ptr] += 1
		} else if code[cnt] == "벤티" {
			ascii[ptr] -= 1
		} else if code[cnt] == "클레" {
			ptr += 1
		} else if code[cnt] == "유라" {
			ptr -= 1
		} else if code[cnt] == "감우" {
			if ascii[ptr] == 0 {
				cnt1 := 1
				for cnt1 != 0 {
					cnt++
					if code[cnt] == "감우" {
						cnt1++
					} else if code[cnt] == "눈나" {
						cnt1--
					} else {
						continue
					}
				}
			}
		} else if code[cnt] == "눈나" {
			if ascii[ptr] != 0 {
				cnt2 := 1
				for cnt2 != 0 {
					cnt--
					if code[cnt] == "눈나" {
						cnt2++
					} else if code[cnt] == "감우" {
						cnt2--
					} else {
						continue
					}
				}
			}
		} else if code[cnt] == "헤으응" {
			output += string(ascii[ptr])
		}
		cnt++
	}
	fmt.Print(output)
	return output
}
