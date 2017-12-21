using Discord.Commands;
using System;
using System.Threading.Tasks;
using Discord;
using System.Collections.Generic;

namespace SlackerBot.Modules
{
    class Cteamup : ModuleBase<SocketCommandContext>
    {
        [Command("TeamUp")]
        public async Task Teamup(IMessage msg)
        {
            List<string> contents = new List<string>(msg.Content.Split(" "));


            string Team1 = "Team 1: ";
            string Team2 = "Team 2: ";
            int factor = 0;
            getFactor(ref factor);
            int spot = 0;
            contents.RemoveAt(0); 

            while(contents.Count >= 0)
            {
                spot = factor % contents.Count;
                Team1 += contents[spot] + " ";
                contents.RemoveAt(spot);
                if (contents.Count >= 0)
                {
                    getFactor(ref factor);
                    spot = factor % contents.Count;
                    Team2 += contents[spot] + " ";
                    contents.RemoveAt(spot);
                }
                await Context.Channel.SendMessageAsync(Team1);
                await Context.Channel.SendMessageAsync(Team2);
            }
        }

        [Command("TeamUpHelp")]
        public async Task Tuphelp()
        {
            await Context.Channel.SendMessageAsync("To use the team up feature: type in the players names. Sepparate Names with a <Space>.");
        }

        private void getFactor(ref int factor)
        {
            factor  += (System.DateTime.Now.Second * System.DateTime.Now.Day) / System.DateTime.Now.Year;
        }
    }
}
