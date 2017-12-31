using Discord.Commands;
using SlackerBot.Modules;
using System;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace SlackerBot
{
    public class TextService
    {
        private readonly IEnumerable<ITextModule<ICommandContext>> _textModules;
        private readonly Dictionary<ProcessTextAttribute, ITextModule> _attributes = new Dictionary<ProcessTextAttribute, ITextModule>();

        public TextService(IEnumerable<ITextModule<ICommandContext>> textModules) {
            _textModules = textModules;
            foreach (var mod in _textModules) {
                var type = mod.GetType();
                foreach (var method in type.GetMethods()) {
                    foreach (var attr in method.GetCustomAttributes(false)) {
                        if (attr is ProcessTextAttribute pattr) {
                            _attributes.Add(pattr, mod);
                        }
                    }
                }
            }
        }

        public async Task<IResult> ExecuteAsync(ICommandContext context) {
            var result = ExecuteResult.FromSuccess();
            foreach (var mod in _textModules) {
                if (mod.RunIf(context.Message.Content)) {
                    mod.Context = context;
                    await mod.Execute();
                }
            }
            return result;
        }
    }
}
