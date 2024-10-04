# Gaudit 🌟

Welcome to **Gaudit**, an elegant and powerful auditing package for Go applications! With Gaudit, you can effortlessly track changes, log activities, and maintain a detailed audit trail of your data.

<p align="center">
    <img src="https://i.ibb.co/51kV1B7/Screenshot-2024-10-05-at-5-56-53-AM.png" alt="Screenshot-2024-10-05-at-5-56-53-AM"/>
</p>

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
