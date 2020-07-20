using Discord;
using Discord.WebSocket;
using System.Threading.Tasks;

namespace SlackerBot
{
    public class Bot
    {
        private const string TOKEN = "Mzg5OTYxNjIyMDUyOTI5NTQ2.DRxCsg.G8C5-hHJRI8kUycc7EfEM8dJOWg";

        private readonly DiscordSocketClient _client;

        public Bot(DiscordSocketClient client) {
            _client = client;
            Setup();
        }

        public async Task StartAsync() {
            await _client.LoginAsync(TokenType.Bot, TOKEN);
            await _client.StartAsync();
            await Task.Delay(-1);
        }

        private void Setup() {

        }
    }
}
