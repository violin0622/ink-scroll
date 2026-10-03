package grpc

import (
	"context"

	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/contents/app"
	"github.com/violin0622/ink-scroll/contents/app/command"
	contv1 "github.com/violin0622/ink-scroll/proto/gen/go/contents/v1"
	ulidpb "github.com/violin0622/ink-scroll/proto/gen/go/ulid"
)

var _ contv1.CmdServiceServer = (*cmdServiceServer)(nil)

func NewCmdService(
	ah app.CommandHandler,
) cmdServiceServer {
	return cmdServiceServer{
		ah: ah,
	}
}

type cmdServiceServer struct {
	contv1.UnimplementedCmdServiceServer

	ah app.CommandHandler
}

func (s *cmdServiceServer) CreateNovel(
	ctx context.Context, req *contv1.CreateNovelRequest) (
	rep *contv1.CreateNovelResponse, err error,
) {
	defer toStatus(&err)

	cmd := command.CreateNovel{
		AuthorID: req.GetAuthorId().Into(),
		Title:    req.GetTitle(),
		Desc:     req.GetDesc(),
		Category: req.GetCategory(),
		Tags:     req.GetTags(),
	}

	id, err := s.ah.CreateNovel(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &contv1.CreateNovelResponse{Id: ulidpb.From(id)}, nil
}

func (s *cmdServiceServer) UpdateNovel(
	ctx context.Context, req *contv1.UpdateNovelRequest) (
	rep *contv1.UpdateNovelResponse, err error,
) {
	defer toStatus(&err)

	cmd := command.UpdateNovelInfo{
		Tags: req.GetTags(),
		Desc: req.GetDesc(),
	}

	if err := s.ah.UpdateNovel(ctx, cmd); err != nil {
		return nil, err
	}

	return &contv1.UpdateNovelResponse{}, nil
}

func (s *cmdServiceServer) UploadChapter(
	ctx context.Context, req *contv1.UploadChapterRequest) (
	rep *contv1.UploadChapterResponse, err error,
) {
	defer toStatus(&err)

	err = s.ah.UploadChapter(ctx, command.UploadChapter{
		NovelID: req.GetNovelId().Into(),
	})
	if err != nil {
		return nil, err
	}

	return &contv1.UploadChapterResponse{Id: ulidpb.From(ulid.Zero)}, nil
}
