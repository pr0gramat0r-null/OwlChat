using Microsoft.Maui;
using Microsoft.Maui.Controls;
using OwlChat.AppShell;

namespace OwlChat;

public partial class App : Application
{
    public App()
    {
        InitializeComponent();
    }

    protected override Microsoft.Maui.Controls.Window CreateWindow(Microsoft.Maui.IActivationState? activationState)
    {
        // Создаём окно, используя AppShell как корневую страницу.
        // Полностью квалифицированное имя используется, чтобы избежать конфликта с пространством имён OwlChat.AppShell.
        return new Microsoft.Maui.Controls.Window(new global::OwlChat.AppShell.AppShell());
    }
}
