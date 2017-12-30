using System;
using System.Collections.Generic;
using System.Text;

namespace SlackerBot.Settings
{
    public interface ISettings
    {
        IEnumerable<string> EnabledChannels { get; }
    }
}
