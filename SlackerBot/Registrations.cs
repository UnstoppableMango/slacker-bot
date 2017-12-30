using Discord.Commands;
using Discord.WebSocket;
using SimpleInjector;
using SlackerBot.Modules;
using SlackerBot.Settings;

namespace SlackerBot
{
    public static class Registrations
    {
        public static void Register(Container container) {
            var assemblies = new[] { typeof(Registrations).Assembly };

            container.RegisterSingleton(() => new CommandService());
            container.RegisterSingleton(() => new DiscordSocketClient());

            container.RegisterCollection(typeof(ITextModule), assemblies);
            container.RegisterSingleton<TextService>();
            container.RegisterSingleton<ISettings, DefaultSettings>();
            container.RegisterSingleton<CommandHandler>();
            container.RegisterSingleton<Bot>();
        }

        public static void RegisterPackages(this Container container) => Register(container);
    }
}
