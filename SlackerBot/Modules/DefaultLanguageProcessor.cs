using Discord.Commands;
using System.Threading.Tasks;

namespace SlackerBot.Modules
{
    public class DefaultLanguageProcessor : ModuleBase<SocketCommandContext>
    {
        public async Task ProcessLanguage(string msg) {
            if (msg.Contains("Jake")
                    || msg.Contains("jake")
                    || msg.Contains("Jacob")
                    || msg.Contains("jacob")) {
                await Context.Channel.SendMessageAsync("It's \"Jakcbo\", friend");
            } else if (msg.Contains("Haley")
                 || msg.Contains("haley")) {
                await Context.Channel.SendMessageAsync("It's \"Helee\", friend");
            } else if (msg.Contains("Kansas")
                 || msg.Contains("kansas")
                 || msg.Contains("Joe")
                 || msg.Contains("joe")) {
                await Context.Channel.SendMessageAsync("It's \"Mr. Angry\", friend");
            } else if (msg.Contains("Chris")
                 || msg.Contains("chris")) {
                await Context.Channel.SendMessageAsync("It's \"Chrus\", friend");
            } else if (msg.Contains("Izaak")
                  || msg.Contains("izaak")) {
                await Context.Channel.SendMessageAsync("It's \"The Greg\", friend");
            } else if (msg.Contains("Stephanie")
                  || msg.Contains("stephanie")
                  || msg.Contains("Steph")
                  || msg.Contains("steph")) {
                await Context.Channel.SendMessageAsync("It's \"freshman\", friend");
            }
        }
    }
}
