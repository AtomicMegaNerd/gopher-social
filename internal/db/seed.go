package db

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/atomicmeganerd/gopher-social/internal/store"
)

var usernames = []string{
	"bob",
	"tina",
	"alice",
	"john",
	"susan",
	"mike",
	"linda",
	"steve",
	"nancy",
	"chris",
	"jane",
	"paul",
	"carol",
	"kevin",
	"laura",
	"brad",
	"karen",
	"gary",
	"betty",
	"eric",
	"diane",
	"brian",
	"heather",
	"doug",
	"pamela",
	"george",
	"rachel",
	"frank",
	"julie",
	"larry",
	"sara",
	"timothy",
	"anna",
	"philip",
	"megan",
	"scott",
	"emma",
	"matthew",
	"amber",
	"jason",
	"rita",
	"zachary",
	"melissa",
	"patrick",
	"kate",
	"victor",
}

var postTitles = []string{
	"10 Tips to Boost Your Coding Productivity",
	"Exploring the Joys of Functional Programming",
	"How to Create Beautiful User Interfaces with Ease",
	"The Ultimate Guide to Debugging Like a Pro",
	"Why Learning New Languages is Fun and Rewarding",
	"Mastering Git: A Comprehensive Tutorial",
	"Building Your First Mobile App: A Step-by-Step Guide",
	"Secrets to Writing Clean and Maintainable Code",
	"The Magic of Automation in Software Development",
	"Embracing Open Source: How to Get Started",
	"Top 5 Tools Every Developer Should Know About",
	"Creating Stunning Visualizations with Data Science",
	"The Exciting Future of Artificial Intelligence",
	"How to Stay Motivated During Long Projects",
	"Collaborating Effectively in Remote Teams",
	"The Benefits of Pair Programming and How to Get Started",
	"A Beginner's Guide to Cloud Computing",
	"How to Balance Work and Life as a Developer",
	"Unlocking Creativity Through Coding Challenges",
	"Celebrating Success: Sharing Your Achievements with the Community",
}

var postContents = []string{
	"Just finished my first open-source contribution, feeling fantastic!",
	"Loving the new features in Go 1.17 - so many improvements!",
	"Started learning Rust today and it's been an amazing journey so far.",
	"Had a productive day debugging and optimizing my code!",
	"Attended an awesome webinar on cloud computing, learned a lot!",
	"My app just hit 1,000 downloads - thank you all for the support!",
	"Pair programming with my colleague was super fun and insightful today!",
	"Experimented with some new UI designs and I’m really happy with the results.",
	"Refactored some legacy code and now it runs twice as fast!",
	"Successfully integrated a new API into our project - seamless experience!",
	"Exploring the world of data science and it’s incredibly fascinating!",
	"Feeling proud of the clean and maintainable code I wrote today.",
	"Collaborated with an amazing team on a challenging project - we nailed it!",
	"Automated some repetitive tasks, saving hours of manual work each week.",
	"Got inspired by a great tech podcast during my morning run.",
	"Finally mastered recursion - can’t wait to use it in my next project!",
	"Deployed my first microservice architecture and everything went smoothly.",
	"Learning about machine learning algorithms has been so rewarding.",
	"Received positive feedback from users, motivating me to keep improving.",
	"Excited to join a hackathon this weekend – ready to innovate and create!",
}

var tagContents = []string{
	"opensource",
	"golang",
	"rust",
	"debugging",
	"cloudcomputing",
	"appdevelopment",
	"pairprogramming",
	"uidesign",
	"codeoptimization",
	"apiintegration",
	"datascience",
	"cleancode",
	"teamcollaboration",
	"automation",
	"techpodcast",
	"recursion",
	"microservices",
	"machinelearning",
	"userfeedback",
	"hackathon",
}

var comments = []string{
	"Great job on that project, your dedication really shows!",
	"Your code is so clean and well-structured, impressive work!",
	"Absolutely love the user interface you designed, it's stunning!",
	"Thanks for your help earlier, it made a big difference!",
	"Fantastic presentation today, very insightful!",
	"Your enthusiasm for learning new things is inspiring!",
	"You’re an awesome team player, always willing to lend a hand.",
	"Keep up the great work, your progress is amazing!",
	"Your debugging skills are top-notch!",
	"Loved reading your latest blog post, very informative!",
	"Congrats on hitting that milestone, well deserved!",
	"You have a real talent for problem-solving!",
	"Appreciate your positive attitude and hard work.",
	"Your innovative ideas always bring fresh perspectives.",
	"Thanks for sharing those resources, super helpful!",
	"You make collaboration so much smoother and enjoyable.",
	"Impressed by how quickly you picked up that new technology!",
	"Your attention to detail is commendable.",
	"You did an excellent job leading that project, kudos!",
	"Your feedback was spot-on and greatly appreciated.",
}

var emailDomains = []string{
	"example.com",
	"mail.com",
	"test.org",
	"demo.net",
	"sample.co",
	"rofl.org",
	"fakesite.com",
	"mydomain.net",
	"email.org",
	"webmail.co",
	"inbox.com",
	"happymail.co",
	"coolmail2000.org",
	"funmail.happy",
}

func Seed(store *store.Storage) {

	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			slog.Error("failed to create user", "error", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			slog.Error("failed to create post", "error", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			slog.Error("failed to create comment", "error", err)
			return
		}
	}

	slog.Info("seeding completed successfully")
}

func generateUsers(n int) []*store.User {
	users := make([]*store.User, n)
	for ix := range n {
		emailDomain := emailDomains[rand.Intn(len(emailDomains))]
		// Generate a random string of 3 digits to append to the username
		suffix := fmt.Sprintf("%d%d%d", rand.Intn(10), rand.Intn(10), rand.Intn(10))
		username := fmt.Sprintf("%s-%s", usernames[rand.Intn(len(usernames))], suffix)
		email := fmt.Sprintf("%s@%s", username, emailDomain)

		password := "123456" // In a real application, ensure passwords are hashed
		users[ix] = &store.User{
			Username: username,
			Email:    email,
			Password: password,
		}
	}
	return users
}

func checkForDuplicates(items []string) {
	seen := make(map[string]struct{})
	for _, item := range items {
		if _, exists := seen[item]; exists {
			panic(fmt.Sprintf("duplicate found: %s", item))
		}
		seen[item] = struct{}{}
	}
}

func generatePosts(n int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, n)
	for ix := range n {
		user := users[rand.Intn(len(users))]
		title := postTitles[rand.Intn(len(postTitles))]
		content := postContents[rand.Intn(len(postContents))]
		tag1 := tagContents[rand.Intn(len(tagContents))]
		tag2 := tagContents[rand.Intn(len(tagContents))]

		posts[ix] = &store.Post{
			UserID:  user.ID,
			Title:   title,
			Content: content,
			Tags:    []string{tag1, tag2},
		}
	}
	return posts
}

func generateComments(n int, users []*store.User, posts []*store.Post) []*store.Comment {
	commentsList := make([]*store.Comment, n)
	for ix := range n {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]
		content := comments[rand.Intn(len(comments))]

		commentsList[ix] = &store.Comment{
			UserID:  user.ID,
			PostID:  post.ID,
			Content: content,
		}
	}
	return commentsList
}
