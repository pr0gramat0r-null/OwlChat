using OwlChat.Models;

namespace OwlChat.Services;

public interface IChatService
{
    IReadOnlyList<ChatRoom> Rooms { get; }
    Task<IReadOnlyList<Message>> GetMessagesAsync(Guid roomId, CancellationToken cancellationToken = default);
    Task<Message> SendMessageAsync(Guid roomId, string sender, string text, CancellationToken cancellationToken = default);
}
