using SimpleInjector;

namespace SlackerBot
{
    public class Program
    {
        static void Main(string[] args) {
            var container = new Container();
            container.RegisterPackages();
            container.Verify();

            var bot = container.GetInstance<Bot>();

            bot.StartAsync().GetAwaiter().GetResult();
        }
    }
}
