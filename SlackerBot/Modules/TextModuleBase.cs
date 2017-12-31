using Discord.Commands;
using System;
using System.Collections.Generic;
using System.Text;
using System.Threading.Tasks;

namespace SlackerBot.Modules
{
    public abstract class TextModuleBase<T> : ITextModule<T>
        where T : ICommandContext
    {
        public T Context { get; set; }

        public abstract bool RunIf(string msg);

        public abstract Task Execute();
    }
}
