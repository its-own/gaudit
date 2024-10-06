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
)

func main() {
    // write your own code
}
```

## 🔧 Configuration

Customize Gaudit to fit your needs. You can configure logging settings, output formats, and more in the `config.go` file.

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

For any questions, suggestions, or feedback, feel free to reach out via [Issues](https://github.com/its-own/gaudit/issues) or contact us directly at [razibulhasan.mithu.com](mailto:razibulhasan.mithu.com).

## 🤝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Thank you for considering **Gaudit** for your auditing needs! We hope it enhances your application's capabilities. Happy coding! 🚀
