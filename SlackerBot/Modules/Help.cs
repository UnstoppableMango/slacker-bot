using Discord.Commands;
using System;
using System.Threading.Tasks;
using Discord;


namespace SlackerBot.Modules
{
    public class Help : ModuleBase<SocketCommandContext>
    {
        string _line_1 = "So you need help huh? Well here are the bassics:";
        string _line_2 = "The basic commands are; !Test !Day !Channel !TeamUp";
        string _line_3 = "Type <??> before a command to get a more detailed description. ??<Command>";

        [Command("Help")]
        public async Task Helpmsg()
        {
            await Context.Channel.SendMessageAsync(_line_1);
            await Context.Channel.SendMessageAsync(_line_2);
            await Context.Channel.SendMessageAsync(_line_3);
        }
    }
}