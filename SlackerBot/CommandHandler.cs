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
                if (msg.Channel.Name.Contains("bot-testing-room")) {
                    var context = new SocketCommandContext(_client, msg);
                    LanguageFilter filter = new LanguageFilter();
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
                    } else if (msg.Content.Contains("Haley")
                         || msg.Content.Contains("haley")) {
                        await context.Channel.SendMessageAsync("It's \"Helee\", friend");
                    } else if (msg.Content.Contains("Kansas")
                         || msg.Content.Contains("kansas")
                         || msg.Content.Contains("Joe")
                         || msg.Content.Contains("joe")) {
                        await context.Channel.SendMessageAsync("It's \"Mr. Angry\", friend");
                    } else if (msg.Content.Contains("Chris")
                         || msg.Content.Contains("chris")) {
                        await context.Channel.SendMessageAsync("It's \"Chrus\", friend");
                    } else if (msg.Content.Contains("Izaak")
                          || msg.Content.Contains("izaak")) {
                        await context.Channel.SendMessageAsync("It's \"The Greg\", friend");
                    } else if (msg.Content.Contains("Stephanie")
                          || msg.Content.Contains("stephanie")
                          || msg.Content.Contains("Steph")
                          || msg.Content.Contains("steph")) {
                        await context.Channel.SendMessageAsync("It's \"freshman\", friend");
                    }
                    if(!filter.CheckMessage(msg))
                    {
                        await context.Message.DeleteAsync(null);
                    }
                }
            }
        }
    }
}
