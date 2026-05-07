namespace OwlChat.Models;

public sealed record Message(Guid Id, Guid ChatId, string Sender, string Text, DateTime SentAtUtc);
