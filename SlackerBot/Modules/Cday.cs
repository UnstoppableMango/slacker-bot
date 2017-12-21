using Discord.Commands;
using System;
using System.Threading.Tasks;

namespace SlackerBot.Modules
{
    public class Cday : ModuleBase<SocketCommandContext>
    {
        [Command("Day")]
        public async Task Daymsg()
        {
            var msg = "";
            var cstZone = TimeZoneInfo.FindSystemTimeZoneById("Central Standard Time");
            switch (TimeZoneInfo.ConvertTimeFromUtc(DateTime.UtcNow, cstZone).DayOfWeek)
            {
                case System.DayOfWeek.Sunday:
                    msg = "https://www.youtube.com/watch?v=iVCSx0S9Upg";
                    break;
                case System.DayOfWeek.Monday:
                    msg = "https://i.pinimg.com/736x/42/bf/ce/42bfce0909b44b4c5dfaa4b48777658a--funny-monday-memes-monday-humor.jpg";
                    break;
                case System.DayOfWeek.Tuesday:
                    msg = "https://img.memecdn.com/tuesday-problems_o_819362.jpg";
                    break;
                case System.DayOfWeek.Wednesday:
                    msg = "https://www.youtube.com/watch?v=du-TY1GUFGk";
                    break;
                case System.DayOfWeek.Thursday:
                    msg = "https://i0.wp.com/wishmeme.com/wp-content/uploads/2017/05/Thursday-Memes-43.jpeg?resize=648%2C648";
                    break;
                case System.DayOfWeek.Friday:
                    msg = "https://i.pinimg.com/736x/d9/4a/88/d94a88cdb796107aac071fae5516b2c9--friday-memes-funny-friday.jpg";
                    break;
                case System.DayOfWeek.Saturday:
                    msg = "https://pics.me.me/another-wild-saturday-night-memes-con-wild-af-theblessedone-25873020.png";
                    break;
                default:
                    msg = "Well... I think it broke...";
                    break;
            }
            await Context.Channel.SendMessageAsync(msg);
        }

        [Command("?Day")]
        public async Task Tuphelp()
        {
            await Context.Channel.SendMessageAsync("Tells you the day mah dude");
        }
    }
}