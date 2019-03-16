package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/weisurya/grpc-golang-playground/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	fmt.Println("Creating blog")
	blog := &blogpb.Blog{
		AuthorId: "Stephane",
		Title:    "My first blog",
		Content:  "Content of the first blog",
	}
	createdBlogResponse, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v\n", createdBlogResponse)

	blogID := createdBlogResponse.GetBlog().GetId()
	fmt.Println("Reading blog")
	_, err = c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "5c8c6cb2af576262fb7056b1",
	})
	if err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}

	readBlogResponse, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: blogID,
	})
	if err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}

	fmt.Printf("Blog was read: %v\n", readBlogResponse)

	// Update blog
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Changed author",
		Title:    "My first updated blog",
		Content:  "Content was updated",
	}

	updateResponse, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: newBlog,
	})
	if err != nil {
		fmt.Printf("Error happened while updating: %v\n", err)
	}

	fmt.Printf("Blog was updated: %v\n", updateResponse)

	// Delete blog
	deleteRes, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: blogID,
	})
	if err != nil {
		fmt.Printf("Error happened while deleting: %v\n", err)
	}

	fmt.Printf("Blog was deleted: %v\n", deleteRes)

	// List blogs
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Error while calling ListBlog RPC: %v\n", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v\n", err)
		}
		fmt.Println(res.GetBlog())
	}
}
