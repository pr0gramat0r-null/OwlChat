# OwlChat (.NET MAUI, C#)

OwlChat is now fully migrated to **.NET MAUI** using **C#**.

## Stack
- .NET MAUI (single-project)
- C# 12
- MVVM-style view models
- In-memory chat/auth services (ready to replace with API)

## Structure
- `src/OwlChat` — MAUI application project.
- `Models` — domain models (user, chat, message).
- `Services` — auth and chat services.
- `ViewModels` — view models for pages.
- `Views` — MAUI pages.

## Run locally
1. Install .NET 9 SDK + MAUI workload.
2. From repository root:
   ```bash
   dotnet build OwlChat.sln
   dotnet run --project src/OwlChat/OwlChat.csproj
   ```

## Notes
The previous Go backend implementation has been removed to align the entire project with the requested MAUI/C# stack.
