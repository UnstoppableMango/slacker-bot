using Discord.Commands;
using System;
using System.Threading.Tasks;
using Discord;
using System.Collections.Generic;

namespace SlackerBot.Modules
{
    public class Cteamup : ModuleBase<SocketCommandContext>
    {
        //Attempts to put listed names onto Teams
        [Command("TeamUp")]
        public async Task Teamup(IMessage msg)
        {
            try
            {
                List<string> contents = new List<string>(msg.Content.Split(" "));


                string Team1 = "Team 1: ";
                string Team2 = "Team 2: ";
                int factor = 0;
                getFactor(ref factor);
                int spot = 0;

                while (contents.Count >= 0)
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
            catch(Exception e)
            {
                await Context.Channel.SendMessageAsync("Looks like something got messed up, Sorry about that. Maybe give 'er another shot!?");
            }
        }


        //The help text for command
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
