using System;
using System.Collections.Generic;
using System.Text;

namespace SlackerBot.Settings
{
    public class DefaultSettings : ISettings
    {
        public IEnumerable<string> EnabledChannels {
            get {
                return new List<string> {
                    "bot-testing-room"
                };
            }
        }
    }
}
