package database

import (
	"context"
	"discord-bot/types"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func CreateMSG(s *discordgo.Session, m *discordgo.MessageCreate, db *mongo.Client, title, data string) {
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "시간zon 에러")
		return
	}

	currentTime := time.Now().In(location).Format("2006-08-15 15:05")

	upload := &types.CreateMSG{
		Title:     title,
		MSG:       data,
		Author:    m.Author.Username,
		CreatedAt: currentTime,
	}

	collection := db.Database("msg").Collection("msg")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	_, err = collection.InsertOne(ctx, upload)
	if err != nil {
		log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "저장 실��")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "완벽하게 숙지 했다.")
}

func GetMSG(db *mongo.Client, word string) ([]types.CreateMSG, bool) {
	if word == "배워" {
		return nil, false
	} else {
		collection := db.Database("msg").Collection("msg")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.D{{"title", word}}
		cur, err := collection.Find(ctx, filter)
		if err != nil {
			log.Println(err)
			return nil, false
		}
		defer cur.Close(ctx)

		var results []types.CreateMSG
		if err := cur.All(ctx, &results); err != nil {
			log.Println(err)
			return nil, false
		}

		return results, true
	}
}

func DeleteMSG(s *discordgo.Session, m *discordgo.MessageCreate) error {
	index := strings.Split(m.Content, " ")
	if len(index) < 3 || index[2] == "" { // 배열 길이 체크 추가
		return errors.New("숫자 안씀")
	}

	count, err := strconv.Atoi(index[2])
	if err != nil {
		return err
	}

	// Discord API의 제한에 따라 최대 100개 메시지를 한 번에 삭제할 수 있습니다.
	if count > 100 {
		count = 100
	}

	messages, err := s.ChannelMessages(m.ChannelID, count, "", "", "")
	if err != nil {
		return err
	}

	// 삭제할 메시지 IDs를 수집합니다.
	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}

	// 메시지를 한 번에 삭제합니다.
	if len(messageIDs) > 0 {
		if err := s.ChannelMessagesBulkDelete(m.ChannelID, messageIDs); err != nil {
			return err
		}
	}

	// 이미지 파일을 전송하기 전 nil 체크 추가
	image, err := os.Open("Untitled.jpg")
	if err != nil {
		return err // 파일 열기 실패 시 에러 반환
	}
	defer image.Close() // 파일을 사용한 후 반드시 닫아줍니다.

	_, err = s.ChannelFileSend(m.ChannelID, "Untitled.jpg", image)
	if err != nil {
		return err // 파일 전송 실패 시 에러 반환
	}

	return nil
}

func SetGamer(s *discordgo.Session, m *discordgo.MessageCreate, db *mongo.Client) {
	userID := m.Author.ID
	upload := &types.Gamer{
		User:   fmt.Sprintf("<@%s>", userID),
		Budget: 1000000,
		Win:    0,
		Games:  0,
	}

	collection := db.Database("msg").Collection("gamer")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	// 사용자 존재 여부 확인
	var existingGamer types.Gamer
	err := collection.FindOne(ctx, bson.M{"user": upload.User}).Decode(&existingGamer)
	if err == nil {
		// 사용자 존재 시
		s.ChannelMessageSend(m.ChannelID, "이미 있는 유저임ㅋ 현실에 리트라이는 없단다ㅋ")
		return
	} else if err != mongo.ErrNoDocuments {
		// 다른 오류 발생 시
		log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "오류가 발생했습니다. 다시 시도해 주세요.")
		return
	}

	// 새로운 사용자 추가
	_, err = collection.InsertOne(ctx, upload)
	if err != nil {
		log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "오류가 발생했습니다. 다시 시도해 주세요.")
		return
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("게이머가 생성되었습니다.", fmt.Sprintf("유저: %s\n자본금: %d원\n승리한 횟수: %d", upload.User, upload.Budget, upload.Win), 16705372))
}

func StartLRGame(s *discordgo.Session, m *discordgo.MessageCreate, db *mongo.Client) {
	userID := m.Author.ID
	collection := db.Database("msg").Collection("gamer")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	var gamer types.Gamer
	err := collection.FindOne(ctx, bson.M{"user": fmt.Sprintf("<@%s>", userID)}).Decode(&gamer)
	if err == mongo.ErrNoDocuments {
		s.ChannelMessageSend(m.ChannelID, "죄송합니다. 초기화되지 않은 유저입니다.")
		return
	} else if err != nil {
		log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "오류가 발생했습니다. 다시 시도해 주세요.")
		return
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced(
		"어서오세요 홀짝 게임랜드입니다.",
		"게임을 시작하려면 배팅을 해주셔야해요.\n\n리챔아 배팅 <돈>\n이기면 배팅한 금액의 2배, 지면 배팅한 금액을 잃게 됩니다.",
		16705372))

	// 베팅 메시지 처리 핸들러 등록
	s.AddHandlerOnce(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID != userID {
			return
		}

		if strings.HasPrefix(m.Content, "리챔아 배팅 ") {
			// 배팅 금액 파싱
			index := strings.Split(m.Content, " ")
			if len(index) < 3 {
				s.ChannelMessageSend(m.ChannelID, "올바른 배팅 금액을 입력하세요.")
				return
			}

			rating, err := strconv.Atoi(index[2])
			if err != nil || rating <= 0 {
				s.ChannelMessageSend(m.ChannelID, "올바른 배팅 금액을 입력하세요.")
				return
			}

			if rating > gamer.Budget {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("돈이 부족합니다. 현재 자본금: %d원", gamer.Budget))
				return
			}

			// 배팅 성공, 홀짝 선택 안내
			s.ChannelMessageSend(m.ChannelID, "홀 또는 짝을 선택해 주세요.")

			// 홀짝 선택을 위한 핸들러 등록 (여기에서 새로운 핸들러로 분리)
			waitForChoice(s, userID, gamer, rating, db)
		}
	})
}

func waitForChoice(s *discordgo.Session, userID string, gamer types.Gamer, rating int, db *mongo.Client) {
	// 홀짝 선택 핸들러 등록
	log.Println("대기중")
	s.AddHandlerOnce(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == userID {
			if m.Content == "홀" || m.Content == "짝" {
				LRGame(s, m, db, gamer, rating)
				return
			} else {
				s.ChannelMessageSend(m.ChannelID, "잘못된 선택입니다. '홀' 또는 '짝'으로 입력해 주세요.")
				// 올바른 선택을 할 때까지 다시 핸들러 대기
				waitForChoice(s, userID, gamer, rating, db)
				return
			}
		}
	})
}

func LRGame(s *discordgo.Session, m *discordgo.MessageCreate, db *mongo.Client, gamer types.Gamer, rating int) {
	gamer.Games += 1
	rand.Seed(time.Now().UnixNano())
	randomOutcome := "홀"
	if rand.Intn(2) == 0 {
		randomOutcome = "짝"
	}

	// 승패 판단
	if m.Content == randomOutcome {
		gamer.Budget += rating // 승리 시 예산 증가
		gamer.Win += 1
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("축하합니다! 당신의 선택은 %s였고, 승리했습니다! 현재 자본금: %d원", randomOutcome, gamer.Budget))
	} else {
		gamer.Budget -= rating // 패배 시 예산 감소
		if gamer.Budget < 0 {
			gamer.Budget = 0 // 예산이 0 이하로 떨어지지 않도록 조정
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("아쉽습니다! 당신의 선택은 %s였고, 패배했습니다. 현재 자본금: %d원", randomOutcome, gamer.Budget))
	}

	// 업데이트된 예산 저장
	collection := db.Database("msg").Collection("gamer")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.UpdateOne(ctx, bson.M{"user": fmt.Sprintf("<@%s>", m.Author.ID)}, bson.M{"$set": bson.M{"budget": gamer.Budget}})
	if err != nil {
		log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "자본금 업데이트 중 오류가 발생했습니다. 다시 시도해 주세요.")
		return
	}
}

func MySTAT(s *discordgo.Session, m *discordgo.MessageCreate, db *mongo.Client) {
	userID := m.Author.ID

	collection := db.Database("msg").Collection("gamer")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	var gamer types.Gamer
	err := collection.FindOne(ctx, bson.M{"user": fmt.Sprintf("<@%s>", userID)}).Decode(&gamer)
	if err == mongo.ErrNoDocuments {
		s.ChannelMessageSend(m.ChannelID, "죄송합니다. 초기화되지 않은 유저입니다.")
		return
	} else if err != nil {
		log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "오류가 발생했습니다. 다시 시도해 주세요.")
		return
	}

	var winRate int
	if gamer.Games > 0 {
		winRate = (gamer.Win * 100) / gamer.Games
	} else {
		winRate = 0 // 경기가 없으면 승률은 0%
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced(
		"내 상황",
		fmt.Sprintf("현재 자본금: %d원\n승률: %d", gamer.Budget, winRate),
		16705372))
}
