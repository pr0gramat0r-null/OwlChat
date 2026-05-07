using Microsoft.Extensions.Logging;
using OwlChat.Services;
using OwlChat.ViewModels;
using OwlChat.Views;

namespace OwlChat;

public static class MauiProgram
{
    public static MauiApp CreateMauiApp()
    {
        var builder = MauiApp.CreateBuilder();
        builder
            .UseMauiApp<App>();

        builder.Services.AddSingleton<IAuthService, AuthService>();
        builder.Services.AddSingleton<IChatService, ChatService>();

        builder.Services.AddSingleton<LoginViewModel>();
        builder.Services.AddSingleton<ChatViewModel>();

        builder.Services.AddSingleton<LoginPage>();
        builder.Services.AddSingleton<ChatPage>();

#if DEBUG
        builder.Logging.AddDebug();
#endif

        return builder.Build();
    }
}
