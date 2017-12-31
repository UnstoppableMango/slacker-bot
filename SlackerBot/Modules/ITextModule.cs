using Discord.Commands;
using System;
using System.Collections.Generic;
using System.Text;
using System.Threading.Tasks;

namespace SlackerBot.Modules
{
    public interface ITextModule
    {
        bool RunIf(string msg);

        Task Execute();
    }

    public interface ITextModule<T> : ITextModule
        where T : ICommandContext
    {
        T Context { get; set; }
    }
}
