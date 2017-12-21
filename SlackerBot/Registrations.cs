using Discord.Commands;
using Discord.WebSocket;
using SimpleInjector;
using SlackerBot.Settings;

namespace SlackerBot
{
    public static class Registrations
    {
        public static void Register(Container container) {
            var assemblies = new[] { typeof(Registrations).Assembly };

            container.RegisterSingleton(() => new CommandService());
            container.RegisterSingleton(() => new DiscordSocketClient());

            container.RegisterSingleton<ISettings, DefaultSettings>();
            container.RegisterSingleton<CommandHandler>();
            container.RegisterSingleton<Bot>();
        }

        public static void RegisterPackages(this Container container) => Register(container);
    }
}
