using System;
using System.Collections.Generic;
using NsqSharp;
using System.Drawing;
using System.Linq;
using System.Threading.Channels;
using System.Text;

namespace Game
{
    internal enum Tile { Empty, Wall, Block, Paddle, Ball }
    
    class Program
    {
        static void Main(string[] args)
        {
            var producer = new Producer("127.0.0.1:4150");
            var consumer = new Consumer("game_output", "game");
            var screen = new Dictionary<Point, Tile>();
            var sig = Channel.CreateBounded<string>(1);
            var h = new GameHandler(screen, sig);
            
            
            consumer.AddHandler(h);
            consumer.ConnectToNsqLookupd("127.0.0.1:4161");

            while (true)
            {
                if (!sig.Reader.TryRead(out var msg)) continue;
                break;
            }

            consumer.Stop();
            Console.WriteLine(h.CountBlocks());
        }
    }

    internal class GameHandler : IHandler
    {
        private readonly Dictionary<Point, Tile> _screen;
        private readonly Channel<string> _ch;
        

        public GameHandler(Dictionary<Point, Tile> screen, Channel<string> ch)
        {
            _screen = screen;
            _ch = ch;
        }

        public void HandleMessage(IMessage message)
        {
            var msg = Encoding.UTF8.GetString(message.Body);
            
            if (msg == "DONE")
            {
                _ch.Writer.TryWrite("DONE");
            } else
            {
                var data = msg.Split(" ").Select(int.Parse).ToArray();
                UpdateState(data);
            }
        }
        private void UpdateState(int[] data)
        {
            var p = new Point(data[0], data[1]);
            _screen[p] = (Tile) data[2];
        }

        public int CountBlocks()
        {
            return _screen.Values.Count(x => x == Tile.Block);
        }

        public void LogFailedMessage(IMessage message)
        {
            throw new NotImplementedException();
        }
    }
}