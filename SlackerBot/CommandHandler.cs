using Discord.Commands;
using Discord.WebSocket;
using SlackerBot.Settings;
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using System.Threading.Tasks;

namespace SlackerBot
{
    public class CommandHandler
    {
        private readonly DiscordSocketClient _client;
        private readonly CommandService _commandService;
        private readonly TextService _textService;
        private readonly ISettings _settings;

        public CommandHandler(
            DiscordSocketClient client,
            CommandService service,
            TextService textService,
            ISettings settings) {
            _client = client;
            _client.MessageReceived += _client_MessageReceived;
            _commandService = service;
            _commandService.AddModulesAsync(Assembly.GetEntryAssembly());
            _textService = textService;
            _settings = settings;
        }

        private async Task _client_MessageReceived(SocketMessage arg) {
            if (!EnabledChannel(arg.Channel.Name)) return;
            if (arg is SocketUserMessage msg) {
                var context = new SocketCommandContext(_client, msg);
                var filter = new LanguageFilter();
                if (!filter.CheckMessage(msg)) {
                    await context.Message.DeleteAsync(null);
                    await context.Channel.SendMessageAsync("No no no, naughty words!");
                    return;
                }

                var argPos = 0;
                if (msg.HasCharPrefix('!', ref argPos)) {
                    //Excecute command
                    var result = await _commandService.ExecuteAsync(context, argPos);
                    if (!result.IsSuccess && result.Error != CommandError.UnknownCommand) {
                        // Failed usage of command, we want to respond
                        await context.Channel.SendMessageAsync(result.ErrorReason);
                    }
                } else if(msg.HasCharPrefix('?', ref argPos))
                {
                    //Look for help text on commands
                    var result = await _commandService.ExecuteAsync(context, argPos);
                    if (!result.IsSuccess && result.Error != CommandError.UnknownCommand)
                    {
                        // Failed usage of command, we want to respond
                        await context.Channel.SendMessageAsync(result.ErrorReason);
                    }
                }

                else {
                    // Fail silently. We don't need the bot responding to EVERY message
                    var result = await _textService.ExecuteAsync(context);
                }
            }
        }

        private bool EnabledChannel(string name) {
            return _settings.EnabledChannels.Contains(name);
        }
    }
}
