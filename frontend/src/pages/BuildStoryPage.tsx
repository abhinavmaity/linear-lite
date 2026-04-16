import { Link } from 'react-router-dom';

type InsightCard = {
  title: string;
  description: string;
};

type Decision = {
  title: string;
  chosen: string;
  avoided: string;
  why: string;
};

type StuckPoint = {
  title: string;
  detail: string;
  unlock: string;
};

type Tone = {
  surface: string;
  edge: string;
  badge: string;
  text: string;
};

const proofPoints = [
  'Focused issue tracker MVP for small engineering teams',
  'React + TypeScript + Vite on the frontend',
  'Go + Gin + PostgreSQL + Redis on the backend',
  'Real frontend/backend contract parity across core workflows',
  // 'Build, smoke, E2E, manual UX, and MVP readiness sign-off',
];

const productSnapshot: InsightCard[] = [
  {
    title: 'What I built',
    description:
      'Linear-lite is a lightweight issue tracking and project planning product built around the workflows teams actually use every day: auth, issues, board/list views, projects, sprints, labels, and dashboard visibility.',
  },
  {
    title: 'What I deliberately did not build',
    description:
      'I kept comments, notifications, attachments, realtime collaboration, advanced analytics, and multi-workspace/RBAC out of scope so the MVP could stay coherent, reliable, and finishable.',
  },
  {
    title: 'Why that matters',
    description:
      'The strongest signal here is not feature count. It is the ability to use AI to accelerate delivery while still making disciplined product and architecture decisions.',
  },
];

const aiApproach: InsightCard[] = [
  {
    title: 'AI helped me move faster on execution',
    description:
      'I used AI to accelerate implementation slices, wiring work, repetitive integrations, and translation from specs into code. That created momentum across the frontend, backend, and validation layers.',
  },
  {
    title: 'AI worked best with strong constraints',
    description:
      'The architecture doc, product scope, milestone artifacts, and validation reports acted as guardrails. That meant AI was not generating in a vacuum; it was operating inside an explicit system of constraints.',
  },
  {
    title: 'Judgment stayed human',
    description:
      'I was still responsible for scope control, trade-offs, contract alignment, acceptance quality, and deciding what was meaningful enough to ship. AI improved throughput, but it did not replace decision-making.',
  },
];

const challenges: InsightCard[] = [
  {
    title: 'Turning a prototype-shaped product into a contract-shaped product',
    description:
      'A polished UI direction is not the same thing as real product parity. One of the biggest challenges was making every major screen behave according to the documented backend contract instead of drifting toward mock-only assumptions.',
  },
  {
    title: 'Keeping the issue workflow consistent across surfaces',
    description:
      'The same issue data had to behave predictably in list view, board view, detail view, archive/restore flows, and dashboard summaries. That consistency work is what makes a product feel trustworthy.',
  },
  {
    title: 'Protecting scope while still shipping something impressive',
    description:
      'AI makes it easy to keep generating. The harder discipline was deciding where the MVP ended and making sure “more” did not quietly become “less coherent.”',
  },
  {
    title: 'Treating validation as product work',
    description:
      'A production-minded build needed more than functioning screens. It needed repeatable build checks, smoke coverage, E2E coverage, manual route acceptance, and a clear readiness decision.',
  },
];

const tradeoffs: Decision[] = [
  {
    title: 'Depth over breadth',
    chosen: 'Shipped a complete core workflow set for issue tracking and planning.',
    avoided: 'Padding the MVP with comments, notifications, attachments, realtime, or enterprise-style controls.',
    why: 'A smaller product that works end to end is more meaningful than a broader product with shallow reliability.',
  },
  {
    title: 'Layered monolith over premature distribution',
    chosen: 'Used a clear backend shape: router -> middleware -> handlers -> services -> repositories -> PostgreSQL.',
    avoided: 'Splitting responsibilities into more complex service boundaries too early.',
    why: 'This kept the system understandable, implementation-friendly, and fast to evolve while staying production-minded.',
  },
  {
    title: 'Real parity over mock convenience',
    chosen: 'Aligned the frontend to real backend responses, validation rules, and supported route behavior.',
    avoided: 'Letting attractive mock patterns imply capabilities the API did not support.',
    why: 'Production credibility comes from contract truth, not from UI optimism.',
  },
  {
    title: 'Validation evidence over demo energy',
    chosen: 'Backed the project with build gates, smoke scripts, E2E coverage, manual UX walkthroughs, and readiness artifacts.',
    avoided: 'Calling it “production-level” because the happy path looked good in a quick demo.',
    why: 'If the goal is real-world impact, the proof has to extend beyond visuals and isolated interactions.',
  },
];

const stuckPoints: StuckPoint[] = [
  {
    title: 'Drawing the MVP boundary honestly',
    detail:
      'The hardest moments were not technical in the narrow sense. They were product decisions about what deserved to exist in the first release and what needed to stay deferred, even if it would have looked impressive on the surface.',
    unlock:
      'I used the documented scope as a forcing function and treated deviations as conscious decisions instead of silent drift.',
  },
  {
    title: 'Reconciling ambition with contract reality',
    detail:
      'Some early UI directions naturally pulled toward richer analytics, writable discussion areas, and expanded interactions. Those ideas were tempting, but they were not supported by the MVP architecture.',
    unlock:
      'I reframed the job as “ship the strongest truthful version of the product,” not “preserve every visually appealing idea.”',
  },
  {
    title: 'Knowing when the product was actually ready',
    detail:
      'There is a big difference between “the features are present” and “the system is ready to be trusted.” I had to decide when to stop building and start proving.',
    unlock:
      'Milestone 5 and Milestone 6 artifacts made the answer concrete: parity first, then hardening, then validation, then sign-off.',
  },
];

const learnings = [
  'AI is most powerful when the problem is already bounded by good docs, clear contracts, and explicit success criteria.',
  'The best use of AI is not replacing engineering judgment. It is increasing the rate at which good judgment can be turned into working software.',
  'Production-minded work is less about generating code quickly and more about keeping code, UX, contracts, validation, and documentation aligned.',
  'Finishing matters. A validated MVP says more about engineering maturity than a long backlog of half-integrated features.',
];

const readinessSignals = [
  'Frontend integrated to real backend APIs for the documented MVP journeys',
  'Backend endpoints implemented for auth, users, projects, sprints, labels, issues, and dashboard',
  'Docker-based full-stack runtime with an explicit migration path',
  'Frontend and backend build gates passing',
  'Smoke coverage for core issue workflow and cache behavior',
  'Critical path browser E2E coverage wired into CI',
  'Manual UX route acceptance completed and MVP readiness signed off',
];

const tones: Tone[] = [
  {
    surface: 'color-mix(in srgb, var(--bg-accent) 14%, var(--bg-elevated))',
    edge: 'var(--bg-accent)',
    badge: 'color-mix(in srgb, var(--bg-accent) 88%, white 12%)',
    text: 'var(--text-primary)',
  },
  {
    surface: 'color-mix(in srgb, var(--bg-secondary) 14%, var(--bg-elevated))',
    edge: 'var(--bg-secondary)',
    badge: 'color-mix(in srgb, var(--bg-secondary) 88%, white 12%)',
    text: 'var(--text-primary)',
  },
  {
    surface: 'color-mix(in srgb, var(--bg-accent-soft) 26%, var(--bg-elevated))',
    edge: 'var(--bg-accent-soft)',
    badge: 'color-mix(in srgb, var(--bg-accent-soft) 86%, white 14%)',
    text: 'var(--text-primary)',
  },
  {
    surface: 'color-mix(in srgb, var(--bg-muted) 76%, var(--bg-elevated))',
    edge: 'var(--border-strong)',
    badge: 'var(--border-strong)',
    text: 'var(--text-primary)',
  },
];

function toneAt(index: number) {
  return tones[index % tones.length];
}

function SectionTitle({ label, title }: { label: string; title: string }) {
  return (
    <div style={{ marginBottom: 16 }}>
      <div className="label" style={{ fontSize: 11, color: 'var(--text-secondary)' }}>
        {label}
      </div>
      <h2 className="label" style={{ fontSize: 30, margin: '8px 0 0 0', lineHeight: 1.1 }}>
        {title}
      </h2>
    </div>
  );
}

function AccentCard({
  title,
  children,
  tone,
}: {
  title: string;
  children: React.ReactNode;
  tone: Tone;
}) {
  return (
    <article
      className="panel-soft"
      style={{
        padding: 16,
        background: tone.surface,
        borderColor: 'var(--border-strong)',
        position: 'relative',
        overflow: 'hidden',
      }}
    >
      <div
        style={{
          width: 52,
          height: 10,
          background: tone.badge,
          border: '2px solid var(--border-strong)',
          boxShadow: 'var(--shadow-soft)',
          marginBottom: 14,
        }}
      />
      <h3 className="label" style={{ margin: 0, fontSize: 16 }}>
        {title}
      </h3>
      <div style={{ marginTop: 10, color: 'var(--text-secondary)', lineHeight: 1.65 }}>{children}</div>
      <div
        style={{
          position: 'absolute',
          right: -22,
          bottom: -22,
          width: 86,
          height: 86,
          background: tone.edge,
          border: '3px solid var(--border-strong)',
          transform: 'rotate(18deg)',
          opacity: 0.18,
        }}
      />
    </article>
  );
}

export function BuildStoryPage() {
  return (
    <main
      className="grid-bg"
      style={{
        minHeight: '100vh',
        padding: '28px 20px 42px',
        background:
          'radial-gradient(circle at 12% 8%, color-mix(in srgb, var(--bg-accent-soft) 36%, transparent) 0, transparent 24%), radial-gradient(circle at 88% 6%, color-mix(in srgb, var(--bg-secondary) 18%, transparent) 0, transparent 20%), radial-gradient(circle at 82% 90%, color-mix(in srgb, var(--bg-accent) 20%, transparent) 0, transparent 22%)',
      }}
    >
      <div style={{ maxWidth: 1160, margin: '0 auto', display: 'grid', gap: 24 }}>
        <header
          className="panel"
          style={{
            padding: 22,
            background:
              'linear-gradient(135deg, color-mix(in srgb, var(--bg-elevated) 82%, var(--bg-accent-soft)) 0%, color-mix(in srgb, var(--bg-elevated) 88%, var(--bg-secondary)) 100%)',
          }}
        >
          <div className="build-story-hero" style={{ display: 'grid', gridTemplateColumns: 'minmax(0, 1.35fr) minmax(280px, 0.95fr)', gap: 18 }}>
            <div>
              <div className="label" style={{ fontSize: 12, color: 'var(--text-secondary)' }}>
                Public Case Study
              </div>
              <h1 className="label" style={{ margin: '8px 0 0 0', fontSize: 'clamp(2.2rem, 4.5vw, 4.4rem)', lineHeight: 0.98 }}>
                USING AI TO SHIP A REAL MVP, NOT JUST A DEMO
              </h1>
              <p style={{ margin: '12px 0 0 0', color: 'var(--text-secondary)', maxWidth: 780, fontSize: 17, lineHeight: 1.6 }}>
                I built Linear-lite to show what responsible AI-assisted software delivery looks like in practice. The point was not to generate as much code
                as possible. The point was to use AI to move faster while staying anchored to architecture, scope, validation, and product truth.
              </p>
              <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap', marginTop: 18 }}>
                <div
                  style={{
                    background: 'var(--bg-accent)',
                    color: 'var(--text-on-accent)',
                    border: '3px solid var(--border-strong)',
                    boxShadow: 'var(--shadow-hard)',
                    padding: '12px 14px',
                    minWidth: 150,
                  }}
                >
                  <div className="label" style={{ fontSize: 10 }}>
                    Build Focus
                  </div>
                  <div style={{ marginTop: 6, fontWeight: 700 }}>Real MVP</div>
                </div>
                <div
                  style={{
                    background: 'var(--bg-secondary)',
                    color: 'var(--text-on-accent)',
                    border: '3px solid var(--border-strong)',
                    boxShadow: 'var(--shadow-hard)',
                    padding: '12px 14px',
                    minWidth: 150,
                  }}
                >
                  <div className="label" style={{ fontSize: 10 }}>
                    Delivery Lens
                  </div>
                  <div style={{ marginTop: 6, fontWeight: 700 }}>AI + Judgment</div>
                </div>
              </div>
            </div>

            <div className="build-story-bento" style={{ display: 'grid', gridTemplateColumns: 'repeat(2, minmax(0, 1fr))', gap: 12, alignContent: 'start' }}>
              <div
                style={{
                  background: 'var(--bg-accent)',
                  color: 'var(--text-on-accent)',
                  border: '3px solid var(--border-strong)',
                  boxShadow: 'var(--shadow-hard)',
                  padding: 16,
                  minHeight: 148,
                }}
              >
                <div className="label" style={{ fontSize: 11 }}>
                  Delivery
                </div>
                <div style={{ marginTop: 10, fontSize: 34, fontWeight: 800 }}>6</div>
                <div style={{ marginTop: 8, lineHeight: 1.45 }}>Milestones to move from implementation planning to MVP readiness.</div>
              </div>
              <div
                style={{
                  background: 'var(--bg-accent-soft)',
                  border: '3px solid var(--border-strong)',
                  boxShadow: 'var(--shadow-hard)',
                  padding: 16,
                  minHeight: 148,
                }}
              >
                <div className="label" style={{ fontSize: 11 }}>
                  Proof
                </div>
                <div style={{ marginTop: 10, fontSize: 34, fontWeight: 800 }}>E2E</div>
                <div style={{ marginTop: 8, lineHeight: 1.45 }}>Validation extended beyond the happy path into repeatable coverage.</div>
              </div>
              <div
                style={{
                  background: 'var(--bg-elevated)',
                  border: '3px solid var(--border-strong)',
                  boxShadow: 'var(--shadow-hard)',
                  padding: 16,
                  minHeight: 148,
                }}
              >
                <div className="label" style={{ fontSize: 11 }}>
                  Architecture
                </div>
                <div style={{ marginTop: 12, fontWeight: 800, fontSize: 20 }}>Layered monolith</div>
                <div style={{ marginTop: 8, color: 'var(--text-secondary)', lineHeight: 1.45 }}>
                  Clear boundaries let AI accelerate execution without causing structural drift.
                </div>
              </div>
              <div
                style={{
                  background: 'color-mix(in srgb, var(--bg-secondary) 20%, var(--bg-elevated))',
                  border: '3px solid var(--border-strong)',
                  boxShadow: 'var(--shadow-hard)',
                  padding: 16,
                  minHeight: 148,
                  display: 'flex',
                  flexDirection: 'column',
                  justifyContent: 'space-between',
                }}
              >
                <div>
                  <div className="label" style={{ fontSize: 11 }}>
                    Product
                  </div>
                  <div style={{ marginTop: 12, fontWeight: 800, fontSize: 20 }}>Open Linear-lite</div>
                </div>
                <Link to="/login" style={{ fontWeight: 700 }}>
                  View product entrypoint
                </Link>
              </div>
            </div>
          </div>
        </header>

        <section
          className="panel"
          style={{
            padding: 22,
            background:
              'linear-gradient(180deg, color-mix(in srgb, var(--bg-accent) 10%, var(--bg-elevated)) 0%, var(--bg-elevated) 100%)',
          }}
        >
          <div style={{ display: 'grid', gridTemplateColumns: 'minmax(0, 1.15fr) minmax(280px, 0.85fr)', gap: 18 }} className="build-story-duo">
            <div>
              <div className="label" style={{ fontSize: 12, color: 'var(--text-secondary)' }}>
                Thesis
              </div>
              <h2 className="label" style={{ fontSize: 'clamp(1.9rem, 3vw, 3.1rem)', margin: '10px 0 0 0', lineHeight: 1 }}>
                AI GAVE ME SPEED. THE WORK STILL NEEDED JUDGMENT.
              </h2>
              <p style={{ marginTop: 14, maxWidth: 980, color: 'var(--text-secondary)', lineHeight: 1.7 }}>
                Linear-lite is a focused issue tracking and planning application for small engineering teams. I used AI throughout the build, but the real
                story is how that speed was shaped by explicit decisions: strong source-of-truth documentation, deliberate MVP boundaries, real integration
                parity, and validation that made the result credible beyond a single polished walkthrough.
              </p>
            </div>
            <div
              className="panel-soft"
              style={{
                padding: 16,
                background: 'color-mix(in srgb, var(--bg-accent-soft) 26%, var(--bg-elevated))',
                display: 'grid',
                alignContent: 'start',
              }}
            >
              <div className="label" style={{ fontSize: 11 }}>
                Core Idea
              </div>
              <p style={{ margin: '10px 0 0 0', lineHeight: 1.6 }}>
                This page is not meant to celebrate AI for its own sake. It is meant to show that I can use AI as a serious execution tool and still keep
                product scope, architecture, and validation under control.
              </p>
            </div>
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: 12, marginTop: 18 }}>
            {proofPoints.map((point) => (
              <AccentCard key={point} title="Signal" tone={toneAt(proofPoints.indexOf(point))}>
                <div style={{ fontWeight: 700, color: 'var(--text-primary)' }}>{point}</div>
              </AccentCard>
            ))}
          </div>
        </section>

        <section
          className="panel"
          style={{
            padding: 22,
            background:
              'linear-gradient(180deg, color-mix(in srgb, var(--bg-secondary) 8%, var(--bg-elevated)) 0%, var(--bg-elevated) 100%)',
          }}
        >
          <SectionTitle label="Product" title="What I Built And Why It Was Worth Building" />
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))', gap: 14 }}>
            {productSnapshot.map((item, index) => (
              <AccentCard key={item.title} title={item.title} tone={toneAt(index + 1)}>
                {item.description}
              </AccentCard>
            ))}
          </div>
        </section>

        <section
          className="panel"
          style={{
            padding: 22,
            background:
              'linear-gradient(180deg, color-mix(in srgb, var(--bg-accent-soft) 18%, var(--bg-elevated)) 0%, var(--bg-elevated) 100%)',
          }}
        >
          <SectionTitle label="AI Workflow" title="How I Actually Used AI During The Build" />
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))', gap: 14 }}>
            {aiApproach.map((item, index) => (
              <AccentCard key={item.title} title={item.title} tone={toneAt(index)}>
                {item.description}
              </AccentCard>
            ))}
          </div>
        </section>

        <section
          className="panel"
          style={{
            padding: 22,
            background:
              'linear-gradient(180deg, color-mix(in srgb, var(--bg-muted) 80%, var(--bg-elevated)) 0%, var(--bg-elevated) 100%)',
          }}
        >
          <SectionTitle label="Challenges" title="What Made This More Than A Simple Build" />
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: 14 }}>
            {challenges.map((item, index) => (
              <AccentCard key={item.title} title={item.title} tone={toneAt(index + 2)}>
                {item.description}
              </AccentCard>
            ))}
          </div>
        </section>

        <section
          className="panel"
          style={{
            padding: 22,
            background:
              'linear-gradient(180deg, color-mix(in srgb, var(--bg-secondary) 9%, var(--bg-elevated)) 0%, color-mix(in srgb, var(--bg-accent) 7%, var(--bg-elevated)) 100%)',
          }}
        >
          <SectionTitle label="Trade-offs" title="Decisions That Kept The Project Honest" />
          <div style={{ display: 'grid', gap: 14 }}>
            {tradeoffs.map((item, index) => {
              const tone = toneAt(index);
              return (
              <article
                key={item.title}
                className="panel-soft"
                style={{
                  padding: 16,
                  background: tone.surface,
                  borderLeft: `10px solid ${tone.edge}`,
                }}
              >
                <h3 className="label" style={{ margin: 0, fontSize: 18 }}>
                  {item.title}
                </h3>
                <p style={{ margin: '10px 0 0 0', lineHeight: 1.6 }}>
                  <strong>Chosen:</strong> {item.chosen}
                </p>
                <p style={{ margin: '8px 0 0 0', lineHeight: 1.6 }}>
                  <strong>Not Chosen:</strong> {item.avoided}
                </p>
                <p style={{ margin: '8px 0 0 0', color: 'var(--text-secondary)', lineHeight: 1.6 }}>
                  <strong>Why:</strong> {item.why}
                </p>
              </article>
            )})}
          </div>
        </section>

        <section
          className="panel"
          style={{
            padding: 22,
            background:
              'linear-gradient(180deg, color-mix(in srgb, var(--bg-accent) 8%, var(--bg-elevated)) 0%, var(--bg-elevated) 100%)',
          }}
        >
          <SectionTitle label="Stuck Points" title="Where I Got Stuck And What Unblocked The Work" />
          <div style={{ display: 'grid', gap: 14 }}>
            {stuckPoints.map((item, index) => {
              const tone = toneAt(index + 1);
              return (
              <article
                key={item.title}
                className="panel-soft"
                style={{
                  padding: 16,
                  background: tone.surface,
                }}
              >
                <h3 className="label" style={{ margin: 0, fontSize: 18 }}>
                  {item.title}
                </h3>
                <p style={{ margin: '10px 0 0 0', lineHeight: 1.65 }}>{item.detail}</p>
                <p style={{ margin: '8px 0 0 0', color: 'var(--text-secondary)', lineHeight: 1.65 }}>
                  <strong>Unlock:</strong> {item.unlock}
                </p>
              </article>
            )})}
          </div>
        </section>

        <section style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: 24 }}>
          <article
            className="panel"
            style={{
              padding: 22,
              background:
                'linear-gradient(180deg, color-mix(in srgb, var(--bg-accent-soft) 16%, var(--bg-elevated)) 0%, var(--bg-elevated) 100%)',
            }}
          >
            <SectionTitle label="Learnings" title="What This Taught Me About AI And Delivery" />
            <div style={{ display: 'grid', gap: 10 }}>
              {learnings.map((item, index) => (
                <div
                  key={item}
                  className="panel-soft"
                  style={{
                    padding: 12,
                    lineHeight: 1.6,
                    background: toneAt(index).surface,
                    borderLeft: `8px solid ${toneAt(index).edge}`,
                  }}
                >
                  {item}
                </div>
              ))}
            </div>
          </article>

          <aside
            className="panel"
            style={{
              padding: 22,
              background:
                'linear-gradient(180deg, color-mix(in srgb, var(--bg-secondary) 14%, var(--bg-elevated)) 0%, color-mix(in srgb, var(--bg-muted) 55%, var(--bg-elevated)) 100%)',
            }}
          >
            <SectionTitle label="Readiness" title="Why This Was Not Left As A Prototype" />
            <div style={{ display: 'grid', gap: 10 }}>
              {readinessSignals.map((item, index) => (
                <div
                  key={item}
                  className="panel-soft"
                  style={{
                    padding: 12,
                    lineHeight: 1.6,
                    background: toneAt(index + 1).surface,
                  }}
                >
                  {item}
                </div>
              ))}
            </div>
          </aside>
        </section>

        <section
          className="panel"
          style={{
            padding: 22,
            background:
              'linear-gradient(135deg, color-mix(in srgb, var(--bg-accent-soft) 42%, var(--bg-elevated)) 0%, color-mix(in srgb, var(--bg-accent) 16%, var(--bg-elevated)) 100%)',
          }}
        >
          <div style={{ display: 'grid', gridTemplateColumns: 'minmax(0, 1.15fr) minmax(260px, 0.85fr)', gap: 18 }} className="build-story-duo">
            <div>
              <div className="label" style={{ fontSize: 11, color: 'var(--text-secondary)' }}>
                Closing
              </div>
              <p style={{ margin: '10px 0 0 0', fontSize: 21, lineHeight: 1.55, maxWidth: 980 }}>
                This project is the clearest way I know how to show my working style: I can leverage AI to build faster, but I do it in a way that stays
                grounded in product scope, engineering discipline, and real validation. That is the kind of AI-assisted software work I want to keep doing.
              </p>
            </div>
            <div
              className="panel-soft"
              style={{
                padding: 16,
                background: 'var(--bg-elevated)',
                display: 'grid',
                alignContent: 'start',
              }}
            >
              <div className="label" style={{ fontSize: 11 }}>
                Takeaway
              </div>
              <div style={{ marginTop: 10, fontWeight: 800, fontSize: 22, lineHeight: 1.3 }}>AI is the accelerator.</div>
              <div style={{ marginTop: 8, color: 'var(--text-secondary)', lineHeight: 1.55 }}>
                Scope, architecture, and validation are what make the output worth trusting.
              </div>
            </div>
          </div>
        </section>
      </div>
    </main>
  );
}
