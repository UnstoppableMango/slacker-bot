using Discord.Commands;
using System;
using System.Collections.Generic;
using System.Text;

namespace SlackerBot.Modules
{
    public abstract class TextModuleBase<T> : ITextModule
        where T : ICommandContext
    {
        protected T Context { get; }
    }
}
