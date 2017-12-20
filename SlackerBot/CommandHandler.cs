using Discord.Commands;
using Discord.WebSocket;
using System.Reflection;
using System.Threading.Tasks;

namespace SlackerBot
{
    public class CommandHandler
    {
        private readonly DiscordSocketClient _client;
        private readonly CommandService _service;

        public CommandHandler(
            DiscordSocketClient client,
            CommandService service) {
            _client = client;
            _client.MessageReceived += _client_MessageReceived;
            _service = service;
            _service.AddModulesAsync(Assembly.GetEntryAssembly());
        }

        private async Task _client_MessageReceived(SocketMessage arg) {
            if (arg is SocketUserMessage msg) {
                var context = new SocketCommandContext(_client, msg);

                var argPos = 0;
                if (msg.HasCharPrefix('!', ref argPos)) {
                    var result = await _service.ExecuteAsync(context, argPos);
                    if (!result.IsSuccess && result.Error != CommandError.UnknownCommand) {
                        await context.Channel.SendMessageAsync(result.ErrorReason);
                    }
                } else if (msg.Content.Contains("Jake")
                    || msg.Content.Contains("jake")
                    || msg.Content.Contains("Jacob")
                    || msg.Content.Contains("jacob")) {
                    await context.Channel.SendMessageAsync("It's \"Jakcbo\", friend");
                }
            }
        }
    }
}
