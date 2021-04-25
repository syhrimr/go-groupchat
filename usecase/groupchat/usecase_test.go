package groupchat

import (
	"errors"
	"reflect"
	"testing"
	"time"

	gom "github.com/golang/mock/gomock"
	"github.com/lolmourne/go-groupchat/model"
	groupchatResource "github.com/lolmourne/go-groupchat/resource/groupchat/mock"
)

func TestUseCase_CreateGroupchat(t *testing.T) {
	controller := gom.NewController(t)
	defer controller.Finish()

	resourceMock := groupchatResource.NewMockDBItf(controller)
	u := UseCase{
		dbRoomRsc: resourceMock,
	}

	type args struct {
		name       string
		adminID    int64
		desc       string
		categoryID int64
	}
	tests := []struct {
		name    string
		args    args
		want    model.Room
		wantErr bool
		mock    func()
	}{
		{
			name: "success create group chat",
			args: args{
				name:       "groupchat nelly",
				adminID:    1,
				desc:       "belajar golang",
				categoryID: 1,
			},
			wantErr: false,
			want: model.Room{
				Name:        "groupchat nelly",
				AdminUserID: 1,
				Description: "belajar golang",
				CategoryID:  1,
			},
			mock: func() {
				resourceMock.EXPECT().CreateRoom(gom.Any(), gom.Any(), gom.Any(), gom.Any()).Return(nil)
			},
		},
		{
			name: "fail create group from db interface",
			args: args{
				name:       "groupchat nelly",
				adminID:    1,
				desc:       "belajar golang",
				categoryID: 1,
			},
			wantErr: true,
			want:    model.Room{},
			mock: func() {
				resourceMock.EXPECT().CreateRoom(gom.Any(), gom.Any(), gom.Any(), gom.Any()).Return(errors.New("DB ERROR"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := u.CreateGroupchat(tt.args.name, tt.args.adminID, tt.args.desc, tt.args.categoryID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCase.CreateGroupchat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UseCase.CreateGroupchat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_GetRoomByID(t *testing.T) {
	controller := gom.NewController(t)
	defer controller.Finish()

	resourceMock := groupchatResource.NewMockDBItf(controller)
	u := UseCase{
		dbRoomRsc: resourceMock,
	}

	type args struct {
		roomID int64
	}
	tests := []struct {
		name    string
		args    args
		want    model.Room
		wantErr bool
		mock    func()
	}{
		{
			name: "success get room",
			args: args{
				roomID: 1,
			},
			want: model.Room{
				RoomID:      1,
				Name:        "Terserah",
				Description: "Deskripsi",
				AdminUserID: 1,
				CategoryID:  1,
				CreatedAt:   time.Date(2020, 2, 1, 1, 1, 1, 1, time.UTC),
			},
			wantErr: false,
			mock: func() {
				resourceMock.EXPECT().GetRoomByID(int64(1)).Return(model.Room{
					RoomID:      1,
					Name:        "Terserah",
					Description: "Deskripsi",
					AdminUserID: 1,
					CategoryID:  1,
					CreatedAt:   time.Date(2020, 2, 1, 1, 1, 1, 1, time.UTC),
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			got, err := u.GetRoomByID(tt.args.roomID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCase.GetRoomByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UseCase.GetRoomByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_multiply(t *testing.T) {
	type args struct {
		x int64
		y int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "success multiply",
			args: args{
				x: 3,
				y: 5,
			},
			want: 15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := multiply(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("multiply() = %v, want %v", got, tt.want)
			}
		})
	}
}
