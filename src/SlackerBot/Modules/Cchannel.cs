﻿using Discord.Commands;
using System;
using System.Threading.Tasks;
using Discord;


namespace SlackerBot.Modules 
{
     public class Cchannel : ModuleBase<SocketCommandContext>
    {

        [Command("Channel")]
        public async Task channel()
        {
            string msg = ("Well around these parts we have 4 channels: " + "\n" 
                + "Otter's: https://www.youtube.com/user/mslisten2urheart " + "\n"
                + "HoboGuy's: https://www.youtube.com/channel/UCBHESdNQuQ4YJEObUiL1T7A" + "\n"
                + "Cyber Slackers': https://www.youtube.com/channel/UCc2-lo0EAmVhYzAy6qnCJig" + "\n"
                + "Chrus': https://www.youtube.com/channel/UCe8JVd1IACixjUGBGEmDxnw");

            await Context.Channel.SendMessageAsync(msg);
        }
    }
} 
 