using OwlChat.Models;

namespace OwlChat.Services;

public interface IAuthService
{
    User? CurrentUser { get; }
    Task<User> SignInAsync(string username, CancellationToken cancellationToken = default);
}
