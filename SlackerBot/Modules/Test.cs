using Discord;
using Discord.Commands;
using System.Threading.Tasks;

namespace SlackerBot.Modules
{
    public class Test : ModuleBase<SocketCommandContext>
    {
        [Command("Test")]
        public async Task TestIt() {
            await Context.Channel.SendMessageAsync("Shit works yo");
        }

        [Command("MessageScope")]
        public async Task TestScope() {
            await Context.Channel.SendMessageAsync(Context.Message.Content.TrimStart('!'));
        }


        [Command("?Test")]
        public async Task Tuphelp()
        {
            await Context.Channel.SendMessageAsync("Tells you if shits workin yo");
        }
    }
}