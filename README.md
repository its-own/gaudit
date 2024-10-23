<p align="center">
    <img src="https://i.postimg.cc/7LK8JWKp/Screenshot-2024-10-07-at-3-27-41-AM.png" height="70%" width="70%" alt="logo">
</p>

## 🏆 Gaudit 

Welcome to **Gaudit**, an elegant and powerful auditing package for Go applications! With Gaudit, you can effortlessly track changes, log activities, and maintain a detailed audit trail of your data.

<p align="center"><a href="#" 
target="_blank"><img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go" alt="go version" /></a>&nbsp;<a href="#" target="_blank"></a></p>

## 📦 Features

- **Comprehensive Auditing**: Automatically log changes for insert and update operations.
- **Flexible Hooks**: Use pre-defined hooks or create your own to customize auditing behavior.
- **Rich Metadata**: Capture essential details such as user actions, timestamps, and IP addresses.
- **Easy Integration**: Seamlessly integrate with your existing Go applications.
- **Extensible**: Extend the package with your own custom functionalities.

## 🚀 Installation

To get started with Gaudit, install it using Go modules:

```bash
go get github.com/its-own/gaudit
```

## 📖 Usage

Here's a quick example of how to use Gaudit in your application:

```go
package main

import (
    "context"
    _ "github.com/its-own/gaudit"
    "github.com/its-own/gaudit/db"
    "github.com/its-own/gaudit/in"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    in.Inject
    ID   primitive.ObjectID `bson:"_id" json:"id"`
    Name string             `bson:"name" json:"name"`
}

func main() {
	// connect to MongoDb 
    client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.MongoUrl).SetRetryWrites(false))
    if err != nil {
        panic(fmt.Sprintf("Failed to create client: %v", err))
    }
    // initialize go audit
    mongo := audit.InitMongo(client, config.MongoDbName, gaudit.New())
    // perform db operation using this client
    performDbOperation(context.Background(), mongo, "test_collection", &User{ID: "123", Name: "John Doe"})
}
// receive interface instead implementation for dependency injection
func performDbOperation(ctx context.Context, c db.NoSql, collection string, obj *User) {
    c.connection.Insert(ctx, collection, obj)
}
```

## 🔧 Configuration

Customize Gaudit to fit your needs. You can configure logging settings, output formats, and more in the `config.go` file.

##  coming soon
```go
// Example of setting configuration
gaudit.SetConfig(gaudit.Config{
    LogLevel: "debug",
    // other configurations...
})
```

## 📚 Documentation

coming soon

## 🧪 Testing

Gaudit comes with a robust set of tests to ensure stability. Run the tests with:

```bash
go test ./...
```

## 🎉 Contributing

We welcome contributions! If you'd like to contribute to Gaudit, please fork the repo and create a pull request. For larger changes, please open an issue first to discuss.

1. Fork the repository
2. Create a new branch (`git checkout -b feature/my-feature`)
3. Make your changes
4. Commit your changes (`git commit -m 'Add some feature'`)
5. Push to the branch (`git push origin feature/my-feature`)
6. Open a Pull Request

## 📧 Contact

For any questions, suggestions, or feedback, feel free to reach out via [Issues](https://github.com/its-own/gaudit/issues) or contact us directly at [razibulhasan.mithu@gmail.com](mailto:razibulhasan.mithu.com).

## 🤝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Thank you for considering **Gaudit** for your auditing needs! We hope it enhances your application's capabilities. Happy coding! 🚀
