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
    }
}
