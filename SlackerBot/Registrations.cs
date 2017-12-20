using Discord.Commands;
using Discord.WebSocket;
using SimpleInjector;
using System;
using System.Collections.Generic;
using System.Text;

namespace SlackerBot
{
    public static class Registrations
    {
        public static void Register(Container container) {
            container.RegisterSingleton(() => new CommandService());
            container.RegisterSingleton(() => new DiscordSocketClient());
            container.RegisterSingleton<CommandHandler>();
            container.RegisterSingleton<Bot>();
        }

        public static void RegisterPackages(this Container container) => Register(container);
    }
}
