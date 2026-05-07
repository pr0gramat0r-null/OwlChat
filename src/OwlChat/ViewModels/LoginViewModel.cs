using System.ComponentModel;
using System.Runtime.CompilerServices;
using System.Windows.Input;
using OwlChat.Services;

namespace OwlChat.ViewModels;

public sealed class LoginViewModel : INotifyPropertyChanged
{
    private readonly IAuthService _authService;
    private string _username = string.Empty;
    private string _status = "Not signed in";

    public event PropertyChangedEventHandler? PropertyChanged;

    public LoginViewModel(IAuthService authService)
    {
        _authService = authService;
        SignInCommand = new Command(async () => await SignInAsync());
    }

    public string Username
    {
        get => _username;
        set
        {
            _username = value;
            OnPropertyChanged();
        }
    }

    public string Status
    {
        get => _status;
        private set
        {
            _status = value;
            OnPropertyChanged();
        }
    }

    public ICommand SignInCommand { get; }

    private async Task SignInAsync()
    {
        var user = await _authService.SignInAsync(Username);
        Status = $"Signed in as {user.Username}";
    }

    private void OnPropertyChanged([CallerMemberName] string? propertyName = null) =>
        PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(propertyName));
}
