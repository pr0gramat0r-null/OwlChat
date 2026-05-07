using OwlChat.Models;

namespace OwlChat.Services;

public sealed class AuthService : IAuthService
{
    public User? CurrentUser { get; private set; }

    public Task<User> SignInAsync(string username, CancellationToken cancellationToken = default)
    {
        var normalized = string.IsNullOrWhiteSpace(username) ? "guest" : username.Trim();
        CurrentUser = new User(Guid.NewGuid(), normalized);
        return Task.FromResult(CurrentUser);
    }
}
