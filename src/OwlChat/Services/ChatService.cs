using OwlChat.Models;

namespace OwlChat.Services;

public sealed class ChatService : IChatService
{
    private readonly List<ChatRoom> _rooms = [new(Guid.NewGuid(), "General")];
    private readonly Dictionary<Guid, List<Message>> _messagesByRoom = new();

    public IReadOnlyList<ChatRoom> Rooms => _rooms;

    public Task<IReadOnlyList<Message>> GetMessagesAsync(Guid roomId, CancellationToken cancellationToken = default)
    {
        _messagesByRoom.TryGetValue(roomId, out var messages);
        return Task.FromResult<IReadOnlyList<Message>>(messages ?? []);
    }

    public Task<Message> SendMessageAsync(Guid roomId, string sender, string text, CancellationToken cancellationToken = default)
    {
        if (!_messagesByRoom.TryGetValue(roomId, out var messages))
        {
            messages = [];
            _messagesByRoom[roomId] = messages;
        }

        var message = new Message(Guid.NewGuid(), roomId, sender, text.Trim(), DateTime.UtcNow);
        messages.Add(message);
        return Task.FromResult(message);
    }
}
