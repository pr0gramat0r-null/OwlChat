using System.Collections.ObjectModel;
using System.ComponentModel;
using System.Runtime.CompilerServices;
using System.Windows.Input;
using OwlChat.Models;
using OwlChat.Services;

namespace OwlChat.ViewModels;

public sealed class ChatViewModel : INotifyPropertyChanged
{
    private readonly IChatService _chatService;
    private readonly IAuthService _authService;
    private string _messageText = string.Empty;

    public event PropertyChangedEventHandler? PropertyChanged;

    public ChatViewModel(IChatService chatService, IAuthService authService)
    {
        _chatService = chatService;
        _authService = authService;
        Messages = [];
        SendCommand = new Command(async () => await SendAsync());
    }

    public ObservableCollection<Message> Messages { get; }

    public string MessageText
    {
        get => _messageText;
        set
        {
            _messageText = value;
            OnPropertyChanged();
        }
    }

    public ICommand SendCommand { get; }

    public async Task InitializeAsync()
    {
        var roomId = _chatService.Rooms[0].Id;
        var existing = await _chatService.GetMessagesAsync(roomId);
        Messages.Clear();
        foreach (var message in existing)
        {
            Messages.Add(message);
        }
    }

    private async Task SendAsync()
    {
        if (string.IsNullOrWhiteSpace(MessageText))
        {
            return;
        }

        var roomId = _chatService.Rooms[0].Id;
        var sender = _authService.CurrentUser?.Username ?? "guest";
        var message = await _chatService.SendMessageAsync(roomId, sender, MessageText);
        Messages.Add(message);
        MessageText = string.Empty;
    }

    private void OnPropertyChanged([CallerMemberName] string? propertyName = null) =>
        PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(propertyName));
}
