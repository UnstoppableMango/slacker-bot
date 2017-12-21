using System;
using System.Collections.Generic;
using System.Text;

using Discord.Commands;
using Discord.WebSocket;
using System.Reflection;
using System.Threading.Tasks;

namespace SlackerBot
{
    class LanguageFilter
    {
        //Checks message for key words, returns true if all clear
        public bool CheckMessage(SocketUserMessage msg)
        {
            if (msg.Content.Contains("asdf"))
                return false;
            return true;
        }
    }
}
