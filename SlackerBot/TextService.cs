using Discord.Commands;
using SlackerBot.Modules;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace SlackerBot
{
    public class TextService
    {
        private readonly IEnumerable<ITextModule> _textModules;
        private readonly Dictionary<ProcessTextAttribute, ITextModule> _attributes = new Dictionary<ProcessTextAttribute, ITextModule>();

        public TextService(IEnumerable<ITextModule> textModules) {
            _textModules = textModules;
            foreach (var mod in _textModules) {
                var type = mod.GetType();
                foreach (var method in type.GetMethods()) {
                    foreach (ProcessTextAttribute attr in method.GetCustomAttributes(false)) {
                        _attributes.Add(attr, mod);
                    }
                }
            }
        }

        public async Task<IResult> ExecuteAsync(SocketCommandContext context) {
            var result = ExecuteResult.FromSuccess();
            return result;
        }
    }
}
