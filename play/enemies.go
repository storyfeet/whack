package play

import (
	"fmt"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"github.com/coderconvoy/engotil"
	"github.com/coderconvoy/engotil/engopoint"
)

type Boxy struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
	engotil.VelocityComponent
	engotil.GCollisionComponent
	ss *common.Spritesheet
}

func NewBoxy(x, y float32, sheet *common.Spritesheet) *Boxy {
	return &Boxy{
		BasicEntity: ecs.NewBasic(),
		SpaceComponent: common.SpaceComponent{
			Position: engo.Point{x, y},
			Width:    40,
			Height:   40,
		},
		RenderComponent: common.RenderComponent{
			Drawable: sheet.Cell(0),
		},
		VelocityComponent: engotil.VelocityComponent{
			Friction: 2,
		},
		GCollisionComponent: engotil.GCollisionComponent{
			Main:  C_ENEMY,
			Group: C_BOY | C_BALL,
			Extra: engo.Point{-20, -20},
		},
		ss: sheet,
	}
}

type BoxSystem struct {
	boxes   []*Boxy
	targets []engotil.Spaceable
}

func (bs *BoxSystem) AddTarget(t engotil.Spaceable) {
	bs.targets = append(bs.targets, t)
}

func (bs *BoxSystem) AddBox(b *Boxy) {
	bs.boxes = append(bs.boxes, b)
}
func (bs *BoxSystem) Update(d float32) {
	if len(bs.targets) == 0 {
		return
	}

	for _, b := range bs.boxes {
		//loop for nearest target
		bcen := engotil.SpaceCenter(*b.GetSpaceComponent())
		nearest := bs.targets[0]
		ncen := engotil.SpaceCenter(*nearest.GetSpaceComponent())
		nd2 := bcen.PointDistanceSquared(ncen)

		for _, t := range bs.targets {
			ncen = engotil.SpaceCenter(*t.GetSpaceComponent())
			td2 := bcen.PointDistanceSquared(ncen)
			if td2 < nd2 {
				nearest = t
				nd2 = td2
			}
		}

		//push towards nearest

		ncen = engotil.SpaceCenter(*nearest.GetSpaceComponent())
		if bcen.X > ncen.X {
			b.Push(engo.Point{-0.1, 0})
		}
		if bcen.X < ncen.X {
			b.Push(engo.Point{0.1, 0})
		}
		if bcen.Y > ncen.Y {
			b.Push(engo.Point{0, -0.1})
		}
		if bcen.Y < ncen.Y {
			b.Push(engo.Point{0, 0.1})
		}

		vc := b.GetVelocityComponent()
		a := engopoint.Angle8(vc.Point)
		rc := b.GetRenderComponent()
		rc.Drawable = b.ss.Cell(a)

	}
}

func (bs *BoxSystem) Remove(e ecs.BasicEntity) {
	bs.targets = engotil.RemoveSpaceable(bs.targets, e)

	bs.boxes = RemoveBoxy(bs.boxes, e)
}

//Killer
type HitSystem struct {
	NPlayers int
}

func (hs *HitSystem) New(w *ecs.World) {
	engo.Mailbox.Listen("GCollisionMessage", func(m engo.Message) {
		cm, ok := m.(engotil.GCollisionMessage)
		if !ok {
			fmt.Println("Could not Convert Message")
			return
		}
		_, isBall := cm.Main.(*Ball)
		_, isBoy := cm.Main.(*Boy)

		_, isBox := cm.Buddy.(*Boxy)
		if isBoy && isBox {
			fmt.Println("Killing")
			engo.SetScene(&MainScene{NPlayers: hs.NPlayers}, true)
		}

		if isBall && isBox {
			fmt.Println("Removing")
			w.RemoveEntity(*cm.Buddy.GetBasicEntity())

		}

	})

}
func (hs *HitSystem) Update(d float32) {}
func (hs *HitSystem) Remove(e ecs.BasicEntity) {
}
