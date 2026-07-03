import '@/App.css'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Slider } from '@/components/ui/slider'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { formatMessage } from '@/utils/format'
import { Sparkles } from 'lucide-react'

function App() {
  return (
    <main className="min-h-screen bg-background px-6 py-16 text-foreground">
      <div className="mx-auto flex max-w-6xl flex-col gap-8">
        <header className="flex flex-wrap items-center justify-between gap-4 rounded-2xl border border-border bg-surface/80 p-6 shadow-sm">
          <div>
            <div className="mb-2 flex items-center gap-2 text-accent">
              <Sparkles className="h-5 w-5" />
              <span className="text-sm font-semibold uppercase tracking-[0.3em]">HabitFlow</span>
            </div>
            <h1 className="text-3xl font-semibold">Design system preview</h1>
            <p className="mt-2 max-w-2xl text-sm text-zinc-400">
              {formatMessage('Tailwind 4, shadcn-style primitives, and the new theme tokens are ready.')}
            </p>
          </div>
          <Badge variant="secondary">Dark mode ready</Badge>
        </header>

        <section className="grid gap-6 lg:grid-cols-[1.25fr_0.75fr]">
          <Card>
            <CardHeader>
              <CardTitle>Core controls</CardTitle>
              <CardDescription>Button, input, dialog, slider, and tabs all render with the new visual language.</CardDescription>
            </CardHeader>
            <CardContent className="flex flex-col gap-4">
              <div className="flex flex-wrap gap-3">
                <Button>Primary action</Button>
                <Button variant="outline">Secondary</Button>
                <Button variant="ghost">Ghost</Button>
              </div>
              <Input placeholder="Search your habits" />
              <Slider defaultValue={[55]} max={100} step={1} />
              <Dialog>
                <DialogTrigger asChild>
                  <Button variant="outline">Open dialog</Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Welcome</DialogTitle>
                    <DialogDescription>Dialog content is styled with the shared theme tokens.</DialogDescription>
                  </DialogHeader>
                </DialogContent>
              </Dialog>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Tabs</CardTitle>
              <CardDescription>Switch between settings and activity views.</CardDescription>
            </CardHeader>
            <CardContent>
              <Tabs defaultValue="overview" className="w-full">
                <TabsList>
                  <TabsTrigger value="overview">Overview</TabsTrigger>
                  <TabsTrigger value="activity">Activity</TabsTrigger>
                </TabsList>
                <TabsContent value="overview" className="mt-4 rounded-xl border border-border bg-background/70 p-4 text-sm text-zinc-300">
                  Daily focus is trending upward with a strong streak.
                </TabsContent>
                <TabsContent value="activity" className="mt-4 rounded-xl border border-border bg-background/70 p-4 text-sm text-zinc-300">
                  Your last three sessions are on track with a 94% completion rate.
                </TabsContent>
              </Tabs>
            </CardContent>
          </Card>
        </section>
      </div>
    </main>
  )
}

export default App
